// simple-edit is a reference monbooru plugin: it rotates images through a
// relay button and crops them on a page of its own behind an open-mode one.
// Every change goes through monbooru's REST API with the token it was issued
// at pairing. 
//

package main

import (
	"bytes"
	"cmp"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"time"

	_ "golang.org/x/image/webp" // monbooru takes webp; the standard library cannot read one
)

const (
	appName = "simple-edit"
	version = "0.1.0"
)

// media = "image" keeps these off videos, animations and comic archives, which
// have no single frame to edit: monbooru renders no button there and drops
// those ids out of a batch.
var buttons = []button{
	{Slot: "detail-actions", Label: "rotate 90", Mode: "relay", Path: "/relay/rotate", Media: "image"},
	{Slot: "batch-bar", Label: "rotate 90", Mode: "relay", Path: "/relay/rotate", Media: "image"},
	{Slot: "detail-actions", Label: "crop", Mode: "open", Path: "/crop?image={image_id}&back={back_url}", Media: "image"},
}

func main() {
	mux := http.NewServeMux()
	// Our own pages reach the browser through monbooru, so they carry the
	// pairing secret like a relay click does.
	mux.HandleFunc("POST /relay/rotate", fromMonbooru(rotateRelay))
	mux.HandleFunc("GET /crop", fromMonbooru(cropPage))
	mux.HandleFunc("GET /crop/preview", fromMonbooru(cropPreview))
	mux.HandleFunc("POST /crop", fromMonbooru(cropSave))
	run(mux)
}

// ---- rotate: a relay button --------------------------------------------

// monbooru gives a relay call 10 s and never retries it, so a batch stops
// short of that and says what it left rather than being cut off mid-image.
const rotateBudget = 8 * time.Second

// rotateRelay answers a click. The payload carries the scope the surface had:
// the one image on a detail page, the selection under the gallery.
func rotateRelay(w http.ResponseWriter, r *http.Request) {
	var click struct {
		ImageIDs []int `json:"image_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&click); err != nil {
		sendJSON(w, http.StatusOK, map[string]any{"ok": false, "message": "could not read the request"})
		return
	}

	deadline := time.Now().Add(rotateBudget)
	var done, failed, left int
	var why error
	for i, id := range click.ImageIDs {
		if i > 0 && time.Now().After(deadline) {
			left = len(click.ImageIDs) - i
			break
		}
		if err := rotate(id); err != nil {
			log.Printf("rotate %d: %v", id, err) // one bad image must not take the batch with it
			failed, why = failed+1, err
			continue
		}
		done++
	}

	message := fmt.Sprintf("rotated %d of %d", done, len(click.ImageIDs))
	if left > 0 {
		message += fmt.Sprintf(", %d left: click again for the rest", left)
	}
	if failed > 0 {
		message += fmt.Sprintf(", %d failed: %v", failed, why)
	}
	// monbooru flashes message; refresh reloads the page so the edit shows.
	sendJSON(w, http.StatusOK, map[string]any{"ok": done > 0, "message": message, "refresh": done > 0})
}

// rotate is the whole contract in three calls: read the bytes, change them,
// write them back. monbooru re-derives the hash, dimensions, thumbnail and
// phash, and keeps the tags, sources and relations.
func rotate(imageID int) error {
	src, format, err := fetchImage(imageID)
	if err != nil {
		return err
	}
	b := src.Bounds()
	turned := image.NewRGBA(image.Rect(0, 0, b.Dy(), b.Dx()))
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			turned.Set(b.Max.Y-1-y, x-b.Min.X, src.At(x, y))
		}
	}
	return replaceImage(imageID, turned, format)
}

// ---- crop: an open-mode button -----------------------------------------

//go:embed crop.html
var cropHTML string

var cropTemplate = template.Must(template.New("crop").Parse(cropHTML))

// cropPage serves what the button opens: a pop-in over the page the user
// was on, served through the mount at /plugins/simple-edit/. That prefix is
// why the page's links are relative.
func cropPage(w http.ResponseWriter, r *http.Request) {
	imageID, err := strconv.Atoi(r.URL.Query().Get("image"))
	if err != nil {
		http.Error(w, "bad image id", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = cropTemplate.Execute(w, map[string]any{
		"ID":   imageID,
		"Back": backURL(r, r.URL.Query().Get("back")),
	})
}

// cropPreview hands the image bytes to the page: the browser holds no token,
// so what it displays comes through us.
func cropPreview(w http.ResponseWriter, r *http.Request) {
	imageID, err := strconv.Atoi(r.URL.Query().Get("image"))
	if err != nil {
		http.Error(w, "bad image id", http.StatusBadRequest)
		return
	}
	raw, err := imageBytes(imageID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	_, _ = w.Write(raw)
}

// cropSave keeps the box the user dragged, over the image's own bytes or
// as a new image linked to it.
func cropSave(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	imageID, err := strconv.Atoi(r.FormValue("image"))
	if err != nil {
		http.Error(w, "bad image id", http.StatusBadRequest)
		return
	}
	src, format, err := fetchImage(imageID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	// The page sends the box in the image's own pixels, so it only has to be
	// clipped to what is there.
	num := func(name string) int { n, _ := strconv.Atoi(r.FormValue(name)); return n }
	box := image.Rect(num("x"), num("y"), num("x")+num("w"), num("y")+num("h")).Intersect(src.Bounds())
	if box.Empty() {
		http.Error(w, "the crop box covers none of the image", http.StatusBadRequest)
		return
	}
	cropped := image.NewRGBA(image.Rect(0, 0, box.Dx(), box.Dy()))
	draw.Draw(cropped, cropped.Bounds(), src, box.Min, draw.Src)

	if r.FormValue("mode") == "version" {
		err = saveAsVersion(imageID, cropped, format)
	} else {
		err = replaceImage(imageID, cropped, format)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	// Leaving the mount closes the pop-in and reloads the page behind it.
	http.Redirect(w, r, backURL(r, r.FormValue("back")), http.StatusSeeOther)
}

// backURL is where the page goes when it is done: what monbooru substituted
// for {back_url}, kept only when it points back at the host the browser reached monbooru at. 
func backURL(r *http.Request, raw string) string {
	host := r.Header.Get("X-Forwarded-Host")
	if u, err := url.Parse(raw); err == nil && u.Host == host {
		return raw
	}
	return cmp.Or(r.Header.Get("X-Forwarded-Proto"), "http") + "://" + host + "/"
}

// ---- what the edits do to the library ----------------------------------

func imageBytes(imageID int) ([]byte, error) {
	raw, _, err := call("GET", api("/images/%d/file", imageID), "", nil)
	return raw, err
}

// fetchImage decodes an image, reporting the format it arrived in so the edit
// can go back the same way.
func fetchImage(imageID int) (image.Image, string, error) {
	raw, err := imageBytes(imageID)
	if err != nil {
		return nil, "", err
	}
	return image.Decode(bytes.NewReader(raw))
}

// encode writes an edited image back out: JPEG stays JPEG, anything else goes
// out as PNG - webp included. monbooru keys the stored type off the bytes and renames the file to
// match, so the extension here only names the upload.
func encode(img image.Image, format string) ([]byte, string, error) {
	var out bytes.Buffer
	if format == "jpeg" {
		err := jpeg.Encode(&out, img, &jpeg.Options{Quality: 95})
		return out.Bytes(), "jpeg", err
	}
	err := png.Encode(&out, img)
	return out.Bytes(), "png", err
}

// upload posts one file as a multipart form, with any plain fields ahead of it.
func upload(url, filename string, body []byte, fields map[string]string) ([]byte, error) {
	var buf bytes.Buffer
	form := multipart.NewWriter(&buf)
	for name, value := range fields {
		if err := form.WriteField(name, value); err != nil {
			return nil, err
		}
	}
	part, err := form.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	if _, err := part.Write(body); err != nil {
		return nil, err
	}
	if err := form.Close(); err != nil {
		return nil, err
	}
	reply, _, err := call("POST", url, form.FormDataContentType(), &buf)
	return reply, err
}

// replaceImage swaps an image's bytes in place. Sending only the file part
// leaves the rest of the row - tags, sources, relations, collections - alone.
func replaceImage(imageID int, img image.Image, format string) error {
	body, ext, err := encode(img, format)
	if err != nil {
		return err
	}
	_, err = upload(api("/images/%d/file", imageID), "edited."+ext, body, nil)
	return err
}

// saveAsVersion pushes the edit as its own image and declares it a newer
// version of the one it came from, so both survive and the relation says which
// came first. via names us as the new row's origin.
func saveAsVersion(parentID int, img image.Image, format string) error {
	body, ext, err := encode(img, format)
	if err != nil {
		return err
	}
	reply, err := upload(api("/images"), "cropped."+ext, body, map[string]string{"via": appName})
	if err != nil {
		return err
	}
	// A push answers with the new image, wrapped when monbooru has something
	// to say about it.
	var created struct {
		ID    int `json:"id"`
		Image struct {
			ID int `json:"id"`
		} `json:"image"`
	}
	if err := json.Unmarshal(reply, &created); err != nil {
		return err
	}
	newID := cmp.Or(created.ID, created.Image.ID)
	if newID == 0 || newID == parentID {
		return nil
	}
	// The image is saved either way, so a link that fails is a log line rather
	// than a failed crop.
	if err := postJSON(api("/relations"), map[string]any{"type": "version", "a": parentID, "b": newID}, nil); err != nil {
		log.Printf("link version %d -> %d: %v", parentID, newID, err)
	}
	return nil
}
