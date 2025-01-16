import { checkJwt } from '/scripts/cookie.js';

document.addEventListener("DOMContentLoaded", async () => {
    console.log(window.location.href);
    if (window.location.pathname == "/") {
        // Beim Laden der Seite prüfen, ob ein gültiger Token vorhanden ist
        /*const token = getCookie("token");
        //console.log("Token gefunden:", token); // Überprüfen, ob ein Token vorhanden ist
      
        if (token) {
          console.log("Token vorhanden, Prüfung beginnt...");
          fetch('/check-token', { method: 'GET', credentials: 'same-origin' })
            .then(response => {
              if (response.ok) {
                console.log("Token ist gültig, Weiterleitung zum Menü...");
                window.location.href = 'menu.html'; // Weiterleitung zum Menü
              } else {
                console.log("Token ist ungültig");
              }
            })
            .catch(error => {
              console.error("Fehler bei der Token-Prüfung", error);
            });
        } else {
            window.location.href = 'login.html'; // Weiterleitung zum Login
        }*/
            if (await checkJwt()) {
                window.location.href = 'menu.html';
            } else {
                window.location.href = 'login.html';
            }
    }
  });
  
  // Funktion zum Abrufen des Cookies
  function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
    return null;
  }
  
document.getElementById('loginForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    const errorMessage = document.getElementById('errorMessage');
    const anmeldebutton = document.getElementById("Anmelden");
    const registrierungsbutton = document.getElementById("Registrieren");

    // Einfache Client-seitige Validierung
    if (username.trim() === '' || password.trim() === '') {
        showError('Bitte füllen Sie alle Felder aus.');
    } else if (username.length < 3) {
        showError('Der Benutzername muss mindestens 3 Zeichen lang sein.');
    } else if (password.length < 3) {
        showError('Das Passwort muss mindestens 3 Zeichen lang sein.');
    } else if (e.submitter == anmeldebutton) {
        try {
            const response = await fetch("/login", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ username, password }),
                credentials: "include",
            });

            if (response.ok) {
                window.location.href = "menu.html";
            } else {
                const error = await response.text();
                showError('Login fehlgeschlagen: ' + error);
            }
        } catch (error) {
            console.log("test>>>");
            console.log(error);
            console.log("<<<");
            if (error.body == "Invalid username or password") {
                showError("Invalid username or password!");
            } else
            {
                showError('Fehler aufgetreten: ' + error.message);
            }
        }
    } else if (e.submitter == registrierungsbutton) {
        try {
            const response = await fetch("/register", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ username, password }),
                credentials: "include",
            });

            if (response.ok) {
                showError('Registrierung wird geprüft', 'success');
            } else {
                const error = await response.text();
                showError('Registrierung fehlgeschlagen: ' + error);
            }
        } catch (error) {
            showError('Fehler aufgetreten: ' + error.message);
        }
    } else {
        showError('Keine Anmeldung oder Registrierung');
    }
});

function showError(message, type = 'error') {
    const errorMessage = document.getElementById('errorMessage');
    errorMessage.textContent = message;
    errorMessage.style.display = 'block';

    if (type === 'success') {
        errorMessage.classList.add('success');
    } else {
        errorMessage.classList.remove('success');
    }
}
