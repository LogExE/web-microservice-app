
const form = document.getElementById("noteForm")

form.onsubmit = postNoteCallback;

async function postNoteCallback(e) {
	e.preventDefault();
	const form = e.currentTarget;

	try {
	    const formData = new FormData(form);
	    const payload = {};
	    for (const [k, v] of formData) {
		payload[k] = v;
	    }

	    const fetchOptions = {
		method: "POST",
		headers: {
		    "Content-Type": "application/json",
		    "Accept": "application/json",
		},
		body: JSON.stringify(payload)
	    };

	    const response = await fetch("/box", fetchOptions);
	    if (!response.ok) {
		const errorMessage = await response.text();
		throw new Error(errorMessage);
	    }
	    
	    console.log(await response.json());

	    window.location.href = "index.html";
	} catch (error) {
	    console.error(error);
	}
};
