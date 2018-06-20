"use strict";

// Helper Functions
function plainAJAX(method, req, data, cb) {
    let httpRequest = new XMLHttpRequest();

    if (cb != null) {
        httpRequest.onreadystatechange = () => {
            if (httpRequest.readyState == 4) {
                // Using promise to stop it blocking the ready state change
                new Promise((resolve, reject) => {

                    if (httpRequest.status != 200) {
                        let eb = document.getElementById("errorBox");
                        eb.innerText = JSON.stringify([method, req, httpRequest.response, httpRequest.status]);
                        return reject(httpRequest.response, httpRequest.status);
                    }

                    cb(httpRequest.response, httpRequest.status);
                    resolve();
                }
                );
            }
        }
            ;
    }

    httpRequest.open(method, req, true);
    httpRequest.send(data);
}

function getAJAX(req, cb) {
    plainAJAX("GET", req, null, cb);
}
function postAJAX(req, data, cb) {
    plainAJAX("POST", req, data, cb);
}

//////

let noteData = {
    dateHead: undefined,
    noteList: undefined,
    newNote: undefined,
    Notes: [],

    setup: function (data, status) {
        this.dateHead = document.getElementById("dateHead");
        this.noteList = document.getElementById("noteList");
        this.newNote = document.getElementById("newNote");

        this.dateHead.innerText = "Today";

        this.Notes = JSON.parse(data);

        this.refreshNoteHTML();

        this.newNote.addEventListener("click", (ev) => {
            this.submitNewNote();
        }
        );
    },

    refreshNoteHTML: function () {
        while (this.noteList.children[0]) {
            this.noteList.removeChild(this.noteList.children[0]);
        }

        for (let i = 0; i < this.Notes.length; i++) {
            this.createNoteLineNode(this.Notes[i]);
        }
    },

    createNoteLineNode: function (line) {
        let lineNode = document.createElement("li");
        lineNode.value = line.ID;

        let noteText = line.Note;

        if (noteText.startsWith(".")) {
            noteText = noteText.substr(1);
            let n = document.createElement("span");
            n.innerText = "⭕"
            n.className = "box";
            lineNode.appendChild(n);
            lineNode.className = "leftpad";

            n.addEventListener("click", (ev) => {
                this.editNote(line.ID, "+" + noteText);
            });
        } else if (noteText.startsWith("+")) {
            noteText = noteText.substr(1);
            let n = document.createElement("span");
            n.innerText = "✔️";
            n.className = "box ticked";
            lineNode.appendChild(n);
            lineNode.className = "leftpad";

            n.addEventListener("click", (ev) => {
                this.editNote(line.ID, "." + noteText);
            });
        } else if (noteText.startsWith("-")) {
            noteText = noteText.substr(1);
            let n = document.createElement("span");
            lineNode.className = "striked";

            lineNode.addEventListener("click", (ev) => {
                this.editNote(line.ID, noteText);
            });
        }

        {
            let n = document.createElement("span");
            n.className = "noteText";
            n.innerText = noteText;
            lineNode.appendChild(n);

            n.addEventListener("click", (ev) => {
                this.editNote(line.ID);
            });
        }


        {
            let n = document.createElement("span");
            let cd = new Date(line.Created);
            let md = new Date(line.Modified);

            if (md) {
                n.className = "modified";
                n.innerText = md.toLocaleTimeString();
            } else {
                n.className = "created";
                n.innerText = cd.toLocaleTimeString();
            }

            lineNode.appendChild(n);

            if (lineNode.className == "striked") {
                n.addEventListener("click", (ev) => {
                    this.editNote(line.ID, noteText);
                });
            } else {
                n.addEventListener("click", (ev) => {
                    this.editNote(line.ID, "-" + noteText);
                });
            }
        }

        this.noteList.appendChild(lineNode);
    },

    editNote: function (id, newVal) {
        const oldline = this.Notes[id];

        let noteVal = newVal;
        while (noteVal == undefined) {
            noteVal = window.prompt("Edit Note", oldline.Note);
        }

        postAJAX("/edit/" + oldline.ID, noteVal, (data, status) => {
            let line = JSON.parse(data);
            this.Notes[line.ID] = line;
            this.refreshNoteHTML();
        });
    },

    submitNewNote: function () {
        let noteVal = window.prompt("New Note");
        postAJAX("/new", noteVal, (data, status) => {
            let line = JSON.parse(data);

            this.Notes.push(line);
            this.createNoteLineNode(line);
        }
        );
    },

    processToday: function (data, status) {
    }
};

// On Load
window.addEventListener("load", (ev) => {
    getAJAX("/today", (data, status) => {
        noteData.setup(data, status);
    }
    );
}
);
