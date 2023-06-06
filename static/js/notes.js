function toggleFormat(command) {
    var sel, range;
    if (window.getSelection) {
        sel = window.getSelection();
        if (sel.rangeCount) {
            range = sel.getRangeAt(0);
            var selectedText = range.toString();
            var newNode = document.createElement(command);
            var nodeText = document.createTextNode(selectedText)
            newNode.appendChild(nodeText)
            //newNode.innerHTML = "<" + command + ">" + selectedText + "</" + command + ">";
            range.deleteContents();
            range.insertNode(newNode);
        }
    } 
}

function toggleList(listType) {
    var sel, range;
    if (window.getSelection) {
        sel = window.getSelection();
        if (sel.rangeCount > 0) {
            range = sel.getRangeAt(0);

            const selectedElements = range.commonAncestorContainer.querySelectorAll('p');
            
            var newList = document.createElement(listType);
            selectedElements.forEach((divElement) => {
                const li = document.createElement('li');
            
                const fragment = document.createDocumentFragment();
                while (divElement.firstChild) {
                    const childNode = divElement.firstChild;
                    fragment.appendChild(childNode);
                }

                li.appendChild(fragment)
                newList.appendChild(li)
            });

            range.deleteContents()
            range.insertNode(newList)
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
        if (command == "insertUnorderedList") { toggleList("ul") };
        if (command == "insertorderedList") { toggleList("ol") };
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
        const newNoteLink = document.createElement("a")
        newNoteLink.href = data['Id']
        const trash_span = document.createElement("span")
        trash_span.setAttribute('onclick',"deleteNote(" + data['Id'] + ")")
        const trash_img = document.createElement("img")
        trash_img.src = "img/trash-24.png"
        trash_img.alt = "delete note"
        trash_span.appendChild(trash_img)

        // Fill single note element with data from response
        //json_data = JSON.parse(data)
        new_note.textContent = data['Title']
        note_id.textContent = data['Id']
        note_id.hidden = true

        // Build element
        newNoteLink.appendChild(new_note)
        newNoteLink.appendChild(note_id)
        note_wrapper.appendChild(newNoteLink)
        note_wrapper.appendChild(trash_span)
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
    const data = JSON.stringify({noteId: parseInt(noteId), content: editorContent});

    fetch('/save', {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: data,
    })
    .then((response) => response.json());
}

// Delete note
function deleteNote(noteId) {
    // TODO: Remove note from note list

    const currentUrl = window.location.href;
    currentNoteId = parseInt(parseNoteIdFromUrl(currentUrl),10);
    const data = JSON.stringify({noteId: currentNoteId, content: ''});
    fetch('/' + noteId + '/delete', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: data,
    })
    .then((response) => response.json())
    .then((data) => {

        if (currentNoteId == noteId){
            window.location.replace('/')
        }

        deleteSingeNoteListItem(noteId)

        console.log('Deletion Success:', data);

        if (data['message'] == 'ok'){
        }
    })
    .catch((error) => {
        console.error('Deletion Error:', error);
    });
}

function deleteSingeNoteListItem(noteId) {
    var parentElements = document.querySelectorAll('.browser-file-item');
    parentElements.forEach(function(parentElement) {
      var childElement = parentElement.querySelector('a[href="' + noteId + '"]');
      if (childElement) {
        parentElement.remove();
      }
    });
  }
  
