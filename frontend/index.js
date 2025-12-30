
const notesListEl = document.getElementById("notesList")
const notesSortSelectorEl = document.getElementById("sortTypeSelector");

async function postLike(id) {
	try {
	    const fetchOptions = {
		method: "POST",
		headers: {
		    "Content-Type": "application/json",
		    "Accept": "application/json",
		},
	    };
	    
	    const response = await fetch(`/like/${id}`, fetchOptions);
	    if (!response.ok) {
		const errorMessage = await response.text();
			throw new Error(errorMessage);
	    }
	    
	    console.log(await response.json());

	    fillNotes();
	} catch (error) {
	    console.error(error);
	}
}

function convertTextToHTMLNote(txt) {
    return txt.split("\n").join("<br>")
}

async function fillNotes() {
    const resp = await fetch("boxes?sortBy=" + notesSortSelectorEl.value);
    const result = await resp.json();
    console.log("recieved object:", result);
    
    // TODO: fix XSS
    notesListEl.innerHTML = result
	.map(note => `<li><blockquote><p>&#xAB;${convertTextToHTMLNote(note.content)}&#xBB;</p></blockquote><p><b>Author</b>: ${note.author}</p><p><b>Likes</b>: ${note.likes}<br/><button onclick="postLike(${note.id})">Like!</button></p></li>`)
	.join("");
}

// TODO: poll in intervals
fillNotes()
	.then(() => console.log("loaded notes, hooray! 💫"))
	.catch(err => console.log("caught error: " + err));

notesSortSelectorEl.onchange = fillNotes;
