package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type CaptionRequest struct {
	GifURL string `json:"gif_url"`
	Text   string `json:"text"`
}

// type CaptionResponse struct {
// 	D string `json:"d"`
// }

func main() {

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/search", searchHandler)
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
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
        }
        #gifGrid {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
            gap: 10px;
            margin-top: 20px;
        }
        .gif-card {
            cursor: pointer;
            border: 2px solid transparent;
            border-radius: 8px;
            overflow: hidden;
            transition: border 0.2s;
        }
        .gif-card:hover {
            border: 2px solid #007bff;
        }
        .gif-card img {
            width: 100%;
            height: auto;
            display: block;
        }
        input[type="text"] {
            padding: 10px;
            font-size: 16px;
            width: 300px;
            max-width: 100%;
            margin: 5px 0;
        }
        button {
            padding: 10px 20px;
            font-size: 16px;
            margin: 5px;
            cursor: pointer;
        }
        #selectedGifPreview {
            max-width: 300px;
            border: 3px solid #007bff;
            border-radius: 8px;
            margin: 10px 0;
        }
    </style>
</head>
<body>
    <h1>Gif Caption Tool</h1>

    <!-- Section 1: Search Form + GIF Grid -->
    <div id="searchSection">
        <form id="searchForm">
            <input type="text" id="query" placeholder="Search (e.g. harry potter)" required><br>
            <button type="submit">Search GIFs</button>
        </form>
        <div id="loading" style="display:none; margin-top:10px;">Searching...</div>
        <div id="gifGrid"></div>
    </div>

    <!-- Section 2: Caption Form (hidden initially) -->
    <div id="captionSection" style="display:none;">
        <h2>Selected GIF</h2>
        <img id="selectedGifPreview" alt="Selected GIF">
        <form id="captionForm">
            <input type="text" id="captionText" placeholder="Caption text" required><br>
            <button type="submit">Generate Caption</button>
            <button type="button" id="backBtn">← Back to Search</button>
        </form>
    </div>

    <!-- Section 3: Result Display (hidden initially) -->
    <div id="resultSection" style="display:none;">
        <h2>Your Captioned GIF</h2>
        <div id="result"></div>
        <div id="actions">
            <button id="shareBtn">Share</button>
            <button id="downloadBtn">Download</button>
            <button id="newBtn">Create Another</button>
        </div>
    </div>

    <script>
    let currentBlob = null;
    let selectedGifURL = null;
    let searchResults = [];

    // Phase 4: Search Implementation
    document.getElementById('searchForm').onsubmit = async (e) => {
        e.preventDefault();

        const query = document.getElementById('query').value;
        document.getElementById('loading').style.display = 'block';
        document.getElementById('gifGrid').innerHTML = '';

        try {
            const response = await fetch('/search?q=' + encodeURIComponent(query));

            if (!response.ok) {
                throw new Error('Search failed');
            }

            const data = await response.json();
            searchResults = data.results;

            displayGifGrid(searchResults);
        } catch (error) {
            alert('Search failed: ' + error.message);
        } finally {
            document.getElementById('loading').style.display = 'none';
        }
    };

    function displayGifGrid(results) {
        const grid = document.getElementById('gifGrid');
        grid.innerHTML = '';

        if (results.length === 0) {
            grid.innerHTML = '<p>No GIFs found. Try a different search.</p>';
            return;
        }

        results.forEach(gif => {
            const gifCard = document.createElement('div');
            gifCard.className = 'gif-card';

            gifCard.innerHTML = '<img src="' + gif.preview_url + '" alt="' + gif.title + '">';
            gifCard.onclick = () => selectGif(gif);

            grid.appendChild(gifCard);
        });
    }

    function selectGif(gif) {
        selectedGifURL = gif.original_url;

        // Show caption section
        document.getElementById('searchSection').style.display = 'none';
        document.getElementById('captionSection').style.display = 'block';

        // Show preview of selected GIF
        document.getElementById('selectedGifPreview').src = gif.preview_url;

        // Focus on text input
        document.getElementById('captionText').focus();
    }

    // Phase 5: Caption Implementation
    document.getElementById('captionForm').onsubmit = async (e) => {
        e.preventDefault();

        document.getElementById('result').innerHTML = 'Generating...';
        document.getElementById('captionSection').style.display = 'none';
        document.getElementById('resultSection').style.display = 'block';

        try {
            const response = await fetch('/caption', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({
                    gif_url: selectedGifURL,
                    text: document.getElementById('captionText').value
                })
            });

            if (!response.ok) {
                throw new Error('Caption generation failed');
            }

            currentBlob = await response.blob();
            const url = URL.createObjectURL(currentBlob);

            document.getElementById('result').innerHTML = '<img src="' + url + '" style="max-width:100%;">';
        } catch (error) {
            alert('Generation failed: ' + error.message);
            document.getElementById('resultSection').style.display = 'none';
            document.getElementById('captionSection').style.display = 'block';
        }
    };

    document.getElementById('backBtn').onclick = () => {
        document.getElementById('captionSection').style.display = 'none';
        document.getElementById('searchSection').style.display = 'block';
        selectedGifURL = null;
    };

    document.getElementById('newBtn').onclick = () => {
        // Reset state
        currentBlob = null;
        selectedGifURL = null;
        searchResults = [];

        // Reset UI
        document.getElementById('result').innerHTML = '';
        document.getElementById('resultSection').style.display = 'none';
        document.getElementById('searchSection').style.display = 'block';
        document.getElementById('query').value = '';
        document.getElementById('captionText').value = '';
        document.getElementById('gifGrid').innerHTML = '';
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

func searchHandler(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
    if !ok || user != "alex" || pass != os.Getenv("AUTH_PASSWORD") {
        w.Header().Set("WWW-Authenticate", `Basic realm="Gophy"`)
        w.WriteHeader(401)
        w.Write([]byte("Unauthorized"))
        return
    }

	query := r.URL.Query().Get("q")
	if query == "" {
		w.WriteHeader(400)
		w.Write([]byte("Query parameter 'q' is required"))
		return
	}

	results := searchGiphy(query)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"results": results,
	})
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
		w.WriteHeader(400)
		w.Write([]byte("Invalid JSON"))
		return
	}

	// Validate gif_url
	if req.GifURL == "" {
		w.WriteHeader(400)
		w.Write([]byte("gif_url is required"))
		return
	}

	// Security check: ensure URL is HTTPS from Giphy
	if !strings.HasPrefix(req.GifURL, "https://") {
		w.WriteHeader(400)
		w.Write([]byte("gif_url must be HTTPS"))
		return
	}

	host := req.GifURL[8:]
	if idx := strings.Index(host, "/"); idx != -1 {
		host = host[:idx]
	}
	if host != "media.giphy.com" && !strings.HasSuffix(host, ".giphy.com") {
		w.WriteHeader(403)
		w.Write([]byte("Only Giphy URLs are allowed"))
		return
	}

	w.Header().Set("Content-Type", "image/gif")
	w.Write(GenerateGif(req.GifURL, req.Text))

	// json.NewEncoder(w).Encode(map[string]string{"short": "hello world"})

}

func GenerateGif(gifURL string, text string) []byte {
	gifBytes := downloadGiphyImage(gifURL)
    captionedGif := addCaption(gifBytes, text)
    return captionedGif
}
