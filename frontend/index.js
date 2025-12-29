
const nlist = document.getElementById("notesList")

async function fillNotes() {
    const resp = await fetch("boxes");
    const result = await resp.json();
    console.log("recieved object:", result);
    
    nlist.innerHTML = result
	.map(note => `<li><p><blockquote>${note.content}</blockquote></p><p>Author: ${note.author}</p><p>Likes: ${note.likes}</p></li>`)
	.join("");
}

// TODO: poll in intervals
fillNotes()
    .then(() => console.log("loaded notes, hooray! 💫"))
    .catch(err => console.log("caught error: " + err));
