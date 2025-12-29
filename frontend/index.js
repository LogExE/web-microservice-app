
const nlist = document.getElementById("notesList")

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

async function fillNotes() {
	const resp = await fetch("boxes");
	const result = await resp.json();
	console.log("recieved object:", result);

	nlist.innerHTML = result
		.map(note => `<li><p><blockquote>${note.content}</blockquote></p><p>Author: ${note.author}</p><p>Likes: ${note.likes} <button onclick="postLike(${note.id})">Like!</button></p></li>`)
		.join("");
}

// TODO: poll in intervals
fillNotes()
	.then(() => console.log("loaded notes, hooray! 💫"))
	.catch(err => console.log("caught error: " + err));
