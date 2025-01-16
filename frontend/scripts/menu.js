import { getCookie, checkJwt } from '/scripts/cookie.js';

/*async function checkjwt() {
    const token = getCookie("token");
        //console.log("Token gefunden:", token); // Überprüfen, ob ein Token vorhanden ist
      
        if (token) {
          console.log("Token vorhanden, Prüfung beginnt...");
          fetch('/welcome', { method: 'GET', credentials: 'same-origin' })
            .then(response => {
              if (response.ok) {
                console.log("Token ist gültig, Weiterleitung zum Menü...");
              } else {
                console.log("Token ist ungültig");
              }
            })
            .catch(error => {
              console.error("Fehler bei der Token-Prüfung", error);
            });
        } else {
            console.log("no Token");
        }
}

function logout() {
    document.cookie = "token=; Max-Age=0; path=/";
    window.location.href = "index.html";
}

/*function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
    return null;
  }*/
    (async () => {
        if (await checkJwt()) {
            console.log("Token valid");
        } else {
            window.location.href = 'login.html';
            console.log("Token not valid");
        }
    })();


function generateCalendar() {
    const today = new Date();
    const currentMonth = today.getMonth();
    const currentYear = today.getFullYear();

    const firstDay = new Date(currentYear, currentMonth, 1);
    const lastDay = new Date(currentYear, currentMonth + 1, 0);

    let date = new Date(firstDay);
    date.setDate(date.getDate() - (date.getDay() + 6) % 7); // Start from Monday

    const calendarBody = document.getElementById('calendarBody');
    
    while (date <= lastDay || date.getDay() !== 1) {
        let week = document.createElement('tr');
        
        for (let i = 0; i < 7; i++) {
            let day = document.createElement('td');
            if (date.getMonth() === currentMonth) {
                day.textContent = date.getDate();
                if (date.toDateString() === today.toDateString()) {
                    day.classList.add('today');
                }
                // Check for appointment
                if (date.getDate() === today.getDate()) {
                    day.classList.add('appointment');
                    day.setAttribute('data-appointment', 'Hello World');
                    day.setAttribute('data-time', '12:00 - 13:00');
                }
            }
            week.appendChild(day);
            date.setDate(date.getDate() + 1);
        }
        
        calendarBody.appendChild(week);
    }

    // Add event listeners for appointments
    const appointmentCells = document.querySelectorAll('.appointment');
    const appointmentDetails = document.getElementById('appointmentDetails');

    appointmentCells.forEach(cell => {
        cell.addEventListener('click', function(e) {
            const rect = this.getBoundingClientRect();
            appointmentDetails.style.display = 'block';
            appointmentDetails.style.top = `${rect.bottom + window.scrollY}px`;
            appointmentDetails.style.left = `${rect.left + window.scrollX}px`;
            appointmentDetails.innerHTML = `
                <h3>${this.getAttribute('data-appointment')}</h3>
                <p>Zeit: ${this.getAttribute('data-time')}</p>
            `;
        });
    });

    // Close appointment details when clicking outside
    document.addEventListener('click', function(e) {
        if (!e.target.classList.contains('appointment') && !appointmentDetails.contains(e.target)) {
            appointmentDetails.style.display = 'none';
        }
    });
}

generateCalendar();

async function fetchWeather() {
    const apiKey = '7785bcd2d45a4595a69222058250101'; // Replace with your API key
    const city = 'Cologne'; // Replace with the city of your choice
    const apiUrl = `https://api.weatherapi.com/v1/current.json?key=${apiKey}&q=${city}&aqi=no`;

    try {
        const response = await fetch(apiUrl);
        const data = await response.json();

        const weatherInfo = `
        <strong>${data.location.name}, ${data.location.country}</strong><br>
        ${data.current.temp_c}°C, ${data.current.condition.text}<br>
        Humidity: ${data.current.humidity}%<br>
        Wind: ${data.current.wind_kph} km/h
        `;

        document.getElementById('weather-info').innerHTML = weatherInfo;
    } catch (error) {
        document.getElementById('weather-info').innerHTML = 'Error fetching weather data';
    }
}

// Fetch the weather on page load
fetchWeather();

async function loadPlugins() {
try {
const response = await fetch('/api/plugins');
if (!response.ok) {
    throw new Error('Fehler beim Laden der Plugins.');
}
const plugins = await response.json();
console.log(plugins);

// Hole den Menücontainer (ul#menu-list)
const menuList = document.getElementById('menu-list');

// Füge für jedes Plugin ein neues Listenelement hinzu
plugins.forEach(plugin => {
    // Erstelle ein neues <li> Element
    const listItem = document.createElement('li');

    // Erstelle einen neuen <a> Tag, der das Plugin repräsentiert
    const button = document.createElement('a');
    button.href = plugin.Path;  // Setze die URL des Plugins
    button.classList.add('button'); // Optional: Füge eine CSS-Klasse hinzu
    button.textContent = plugin.Name; // Setze den Text des Links

    // Füge den Button als Kind des Listenelements hinzu
    listItem.appendChild(button);

    // Hänge das Listenelement an die Menü-Liste an
    menuList.appendChild(listItem);
});
} catch (error) {
console.error('Fehler:', error);
}
}

// Lade die Plugins, wenn das Dokument vollständig geladen ist
document.addEventListener('DOMContentLoaded', loadPlugins);