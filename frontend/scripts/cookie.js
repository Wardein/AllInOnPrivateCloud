function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
    return null;
  }

  async function checkJwt() {
    try {
        const response = await fetch('/check-token', {
            method: 'GET',
            credentials: 'same-origin', // TODO: include?
        });

        if (response.ok) {
            const data = await response.text();
            console.log("Erfolg:", data);
            return true;
        } else {
            console.error("Ungültiges Token:", response.status);
            return false;
        }
    } catch (error) {
        console.error("Fehler beim Überprüfen des Tokens:", error);
        return false;
    }
}

export {getCookie , checkJwt};