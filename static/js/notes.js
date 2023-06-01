function toggleFormat(command) {
    var sel, range;
    if (window.getSelection) {
        sel = window.getSelection();
        if (sel.rangeCount) {
            range = sel.getRangeAt(0);
            var selectedText = range.toString();
            var newNode = document.createElement("span");
            newNode.innerHTML = "<" + command + ">" + selectedText + "</" + command + ">";
            range.deleteContents();
            range.insertNode(newNode);
        }
    } 
}

/*

document.getElementById('boldBtn').addEventListener('click', function() {
    toggleFormat("b");
});

document.getElementById('italicBtn').addEventListener('click', function() {
    toggleFormat("i");
});*/

const elements = document.querySelectorAll(".btn")
elements.forEach(element => {
    element.addEventListener("click", () => {
        let command = element.dataset['element']
        if (command == "bold") { toggleFormat("b") };
        if (command == "italic") { toggleFormat("i") };
        if (command == "underline") { toggleFormat("u") };
    })
})


/* Server functionality below */

// Save note on keypress
document.addEventListener('keydown', function(event) {
    if (event.ctrlKey && event.key === 's') {
        saveNoteState();
        }
    });

// Save current state of note every 5 sec
setInterval(() => {saveNoteState();}, 5000); // Save every 5 seconds

document.getElementById('add-note').addEventListener('click', function() {
    createNewNote();
});

function createNewNote() {
    // Fetch new note
    fetch('/new_note', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        }
    })
    .then((response) => response.json())
    .then((data) => {

        // Create new single note element for notes list
        const note_wrapper = document.createElement('div')
        note_wrapper.classList.add('browser-file-item')
        const new_note = document.createElement('span')
        new_note.classList.add('note-name')
        const note_id = document.createElement('span')
        note_id.classList.add('note-id')

        /* TODO: add trashbin and new note link */
        const tash_span = document.createElement("span")

        // Fill single note element with data from response
        //json_data = JSON.parse(data)
        new_note.textContent = data['Title']
        note_id.textContent = data['Id']
        note_id.hidden = true

        // Build element
        note_wrapper.appendChild(new_note)
        note_wrapper.appendChild(note_id)
        notes_list = document.getElementsByClassName('notes-list')

        // Insert on top of notes list
        notes_list[0].insertBefore(note_wrapper,notes_list[0].firstChild)
        
    })
}

function parseNoteIdFromUrl(url) {
    const regex = /\/([^\/]+)$/;
    const match = url.match(regex);
    
    if (match && match.length > 1) {
      return match[1];
    }
    
    return null;
  }

function saveNoteState() {
    const editorContent = document.getElementById('editor').innerHTML;
    const currentUrl = window.location.href;
    const noteId = parseNoteIdFromUrl(currentUrl)
    const data = JSON.stringify({noteId: noteId, content: editorContent});

    fetch('/save', {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: data,
    })
    .then((response) => response.json())
    .then((data) => {
        console.log('Success:', data);
    })
    .catch((error) => {
        console.error('Error:', error);
    });
}

// Delete note
function deleteNote() {
    // TODO: Remove note from note list

    const currentUrl = window.location.href;
    noteId = parseNoteIdFromUrl(currentUrl);
    const data = JSON.stringify({noteId: noteId, content: ''});
    fetch('/delete', {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
        },
        body: data,
    })
    .then((response) => response.json())
    .then((data) => {
        console.log('Success:', data);
    })
    .catch((error) => {
        console.error('Error:', error);
    });
}
