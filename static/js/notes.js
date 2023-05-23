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

document.getElementById('boldBtn').addEventListener('click', function() {
    toggleFormat("b");
});

document.getElementById('italicBtn').addEventListener('click', function() {
    toggleFormat("i");
});

setInterval(() => {
    const editorContent = document.getElementById('editor').innerHTML;
    const data = JSON.stringify({content: editorContent});

    fetch('/save', {
        method: 'POST',
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
}, 5000); // Save every 5 seconds

document.getElementById('add-note').addEventListener('click', function() {
    createNewNote();
});

function createNewNote() {
    fetch('/new_note', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        }
    })
    .then((response) => response.json())
    .then((data) => {
        const note_wrapper = document.createElement('div')
        note_wrapper.classList.add('browser-file-item')
        const new_note = document.createElement('span')
        new_note.classList.add('note-name')
        const node_id = document.createElement('span')
        note_id.classList.add('note-id')

        new_note.textContent = 'Untitled'
        note_id.textContent = data['noteId']
        note_id.hidden = true

        note_wrapper.appendChild(new_note)
        note_wrapper.appendChild(note_id)
        notes_list = document.getElementsByClassName('notes-list')
        notes_list.prepend(note_wrapper)
    })
}