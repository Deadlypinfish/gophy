package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type CaptionRequest struct {
	Query string `json:"query"`
	Text string `json:"text"`
}

// type CaptionResponse struct {
// 	D string `json:"d"`
// }

func main() {

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/caption", caption)

	fmt.Println("Listening on :8080")
	http.ListenAndServe("127.0.0.1:8080", nil)

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
    if !ok || user != "alex" || pass != os.Getenv("AUTH_PASSWORD") {
        w.Header().Set("WWW-Authenticate", `Basic realm="Gophy"`)
        w.WriteHeader(401)
        w.Write([]byte("Unauthorized"))
        return
    }

    html := `
<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body>
    <h1>Gif Caption Tool</h1>
    <form id="gifForm">
        <input type="text" id="query" placeholder="Search (e.g. harry potter)" required><br>
        <input type="text" id="text" placeholder="Caption text" required><br>
        <button type="submit">Generate</button>
    </form>
    <div id="result"></div>
    <div id="actions" style="display:none;">
        <button id="shareBtn">Share</button>
        <button id="downloadBtn">Download</button>
    </div>
    
    <script>
    let currentBlob = null;
    
    document.getElementById('gifForm').onsubmit = async (e) => {
        e.preventDefault();
        document.getElementById('result').innerHTML = 'Generating...';
        
        const response = await fetch('/caption', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({
                query: document.getElementById('query').value,
                text: document.getElementById('text').value
            })
        });
        
        currentBlob = await response.blob();
        const url = URL.createObjectURL(currentBlob);
        
        document.getElementById('result').innerHTML = '<img src="' + url + '" style="max-width:100%;">';
        document.getElementById('actions').style.display = 'block';
    };
    
    document.getElementById('shareBtn').onclick = async () => {
        if (!currentBlob) return;
        
        const file = new File([currentBlob], 'caption.gif', { type: 'image/gif' });
        
        if (navigator.share && navigator.canShare && navigator.canShare({files: [file]})) {
            try {
                await navigator.share({
                    files: [file],
                    title: 'Custom Gif'
                });
            } catch (err) {
                if (err.name !== 'AbortError') {
                    alert('Share failed: ' + err.message);
                }
            }
        } else {
            alert('Sharing not supported on this device - use Download instead');
        }
    };
    
    document.getElementById('downloadBtn').onclick = () => {
        if (!currentBlob) return;
        
        const url = URL.createObjectURL(currentBlob);
        const a = document.createElement('a');
        a.href = url;
        a.download = 'caption.gif';
        a.click();
    };
    </script>
</body>
</html>
    `
    w.Header().Set("Content-Type", "text/html")
    w.Write([]byte(html))
}

func caption(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
    if !ok || user != "alex" || pass != os.Getenv("AUTH_PASSWORD") {
        w.Header().Set("WWW-Authenticate", `Basic realm="Gophy"`)
        w.WriteHeader(401)
        w.Write([]byte("Unauthorized"))
        return
    }

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return
	}

	var req CaptionRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		fmt.Println("Error parsing json:", err)
		return
	}

	// var res CaptionResponse
	// res.D = fmt.Sprintf("q:%s, t:%s\n", req.Query, req.Text)


	w.Header().Set("Content-Type", "image/gif")
	w.Write(GenerateGif(req.Query, req.Text))

	// json.NewEncoder(w).Encode(map[string]string{"short": "hello world"})

}

func GenerateGif(q string, t string) []byte {
	gifBytes := getGiphyImage(q)
    captionedGif := addCaption(gifBytes, t)
    return captionedGif
}
