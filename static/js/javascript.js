// Fonction pour générer une couleur aléatoire
function getRandomColor() {
    const letters = '0123456789ABCDEF';
    let color = '#';
    for (let i = 0; i < 6; i++) {
        color += letters[Math.floor(Math.random() * 16)];
    }
    return color;
}

// Appliquer un dégradé de couleur aléatoire au fond de la page
function setBackground() {
    const color1 = getRandomColor();
    const color2 = getRandomColor();
    const color3 = getRandomColor();

    // Appliquer le dégradé uniquement au fond de la page
    document.body.style.background = `linear-gradient(45deg, ${color1}, ${color2}, ${color3})`;

    // Appliquer les couleurs aux boutons sans affecter le texte des autres éléments
    setButtonAnimations(color1, color2, color3);
}

// Appliquer l'animation des boutons avec des couleurs dynamiques
function setButtonAnimations(color1, color2, color3) {
    const buttons = document.querySelectorAll(".button");

    buttons.forEach(button => {
        button.style.transition = "background-color 0.3s ease, color 0.3s ease, box-shadow 0.3s ease"; // Animation des couleurs

        button.addEventListener("mouseenter", () => {
            button.style.backgroundColor = "#FF3366";  // Fond spécifique au bouton
            button.style.color = "white";  // Texte en blanc
            button.style.boxShadow = `0px 8px 20px rgba(255, 51, 102, 0.5)`; // Ombre correspondant à la couleur des boutons
        });

        button.addEventListener("mouseleave", () => {
            button.style.backgroundColor = "white";  // Fond blanc par défaut
            button.style.color = "#FF3366";  // Texte en couleur spécifique
            button.style.boxShadow = "0px 5px 10px rgba(0, 0, 0, 0.2)";
        });
    });
}

// Appliquer les changements de couleur de fond et d'animation des boutons au chargement de la page
document.addEventListener("DOMContentLoaded", () => {
    setBackground();  // Applique les couleurs au chargement de la page
});
