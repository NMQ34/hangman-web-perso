package main

import (
    "html/template"
    "log"
    "math/rand"
    "net/http"
    "os"
    "strings"
    "time"
)

var (
    targetWord      string
    guessedLetters  []string
    attempts        int
    maxAttempts     int
    difficulty      string
    wordsFile       string
    username        string
)

func main() {
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/start", startHandler)
    http.HandleFunc("/hangman", hangmanHandler)
    http.HandleFunc("/reset", resetHandler)
	http.HandleFunc("/scoreboard", scoreboardHandler)

    // Serve static files (CSS, JS)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../static"))))

    log.Println("Server started at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html") // Définir le Content-Type

    tmpl, err := template.ParseFiles("../templates/index.html")
    if err != nil {
        http.Error(w, "Error loading template", http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, nil)
}

func startHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        r.ParseForm()
        username = r.FormValue("username")
        difficulty = r.FormValue("difficulty")

        // Choisir le fichier selon la difficulté
        switch difficulty {
        case "EASY":
            wordsFile = "../words/words.txt"
        case "NORMAL":
            wordsFile = "../words/words2.txt"
        case "HARD":
            wordsFile = "../words/words3.txt"
        case "BONUS":
            wordsFile = "../words/words_bonus.txt"
        default:
            wordsFile = "../words/words.txt"
        }

        targetWord = getRandomWord(wordsFile)
        guessedLetters = []string{}
        attempts = 0
        maxAttempts = 10

        http.Redirect(w, r, "/hangman", http.StatusSeeOther)
    } else {
        http.Redirect(w, r, "/", http.StatusSeeOther)
    }
}

func getRandomWord(file string) string {
    words, err := os.ReadFile(file)
    if err != nil {
        log.Fatal(err)
    }
    wordList := strings.Split(string(words), "\n")
    rand.Seed(time.Now().Unix())
    return strings.ToUpper(strings.TrimSpace(wordList[rand.Intn(len(wordList))]))
}

func hangmanHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    if r.Method == http.MethodPost {
        r.ParseForm()
        input := strings.ToUpper(r.FormValue("letter"))

        // Gestion de l'entrée : lettre ou mot complet
        if len(input) > 1 {
            // Deviner un mot entier
            if input == targetWord {
                guessedLetters = strings.Split(targetWord, "") // Marquer toutes les lettres devinées
            } else {
                attempts += 2 // Perdre 2 vies pour un mot incorrect
                if attempts > maxAttempts {
                    attempts = maxAttempts // S'assurer que le joueur ne dépasse pas la limite
                }
            }            
        } else if len(input) == 1 {
            // Deviner une lettre unique
            if !strings.Contains(strings.Join(guessedLetters, ""), input) {
                guessedLetters = append(guessedLetters, input)
                if !strings.Contains(targetWord, input) {
                    attempts++
                }
            }
        }
    }

    // Vérifiez si le jeu est terminé
    if attempts >= maxAttempts || maskWord(targetWord, guessedLetters) == targetWord {
        updateScoreboard(username, maskWord(targetWord, guessedLetters) == targetWord)

        tmpl, err := template.ParseFiles("../templates/result.html")
        if err != nil {
            log.Printf("Error loading result template: %v", err)
            http.Error(w, "Error loading result page", http.StatusInternalServerError)
            return
        }

        tmpl.Execute(w, struct {
            Won      bool
            Word     string
            Username string
        }{
            Won:      maskWord(targetWord, guessedLetters) == targetWord,
            Word:     targetWord,
            Username: username,
        })
        return
    }

    tmpl, err := template.ParseFiles("../templates/hangman.html")
    if err != nil {
        http.Error(w, "Error loading template", http.StatusInternalServerError)
        return
    }

    tmpl.Execute(w, struct {
        Word         string
        Guessed      []string
        AttemptsLeft int
        Stickman     string
        GameOver     bool
        Won          bool
    }{
        Word:         maskWord(targetWord, guessedLetters),
        Guessed:      guessedLetters,
        AttemptsLeft: maxAttempts - attempts,
        Stickman:     renderStickman(attempts),
        GameOver:     attempts >= maxAttempts,
        Won:          maskWord(targetWord, guessedLetters) == targetWord,
    })
}

func maskWord(word string, guessed []string) string {
    masked := ""
    for _, c := range word {
        letter := string(c)
        if contains(guessed, letter) {
            masked += letter
        } else {
            masked += " _ "
        }
    }
    return masked
}

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
    // Réinitialise les variables globales
    guessedLetters = []string{}
    attempts = 0
    targetWord = getRandomWord(wordsFile)

    http.Redirect(w, r, "/hangman", http.StatusSeeOther)
}

func renderStickman(attempts int) string {
    // Représentations du pendu selon les étapes
    stickmanStages := []string{
        `
         
            
             
             
             
             
        =========`,
        `
         
             |
             |
             |
             |
             |
        =========`,
        `
         -----
             |
             |
             |
             |
             |
        =========`,
        `
         -----
         |   |
             |
             |
             |
             |
        =========`,
        `
         -----
         |   |
         O   |
             |
             |
             |
        =========`,
        `
         -----
         |   |
         O   |
         |   |
             |
             |
        =========`,
        `
         -----
         |   |
         O   |
        /|   |
             |
             |
        =========`,
        `
         -----
         |   |
         O   |
        /|\  |
             |
             |
        =========`,
        `
         -----
         |   |
         O   |
        /|\  |
        /    |
             |
        =========`,
        `
         -----
         |   |
         O   |
        /|\  |
        / \  |
             |
        =========`,
    }

    if attempts < len(stickmanStages) {
        return stickmanStages[attempts]
    }
    return stickmanStages[len(stickmanStages)-1] // Dernière étape (pendu complet)
}

type Score struct {
    Rank     int
    Username string
    Score    int
}

var scoreboard []Score

func scoreboardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html") 
    tmpl, err := template.ParseFiles("../templates/scoreboard.html")
    if err != nil {
        log.Printf("Error loading scoreboard template: %v", err)
        http.Error(w, "Error loading scoreboard page", http.StatusInternalServerError)
        return
    }

    tmpl.Execute(w, struct {
        Scores []Score
    }{
        Scores: scoreboard,
    })
}

func updateScoreboard(username string, won bool) {
    score := 0
    if won {
        score = 10 - attempts // Exemple : 10 points moins les tentatives
    }
    scoreboard = append(scoreboard, Score{
        Rank:     len(scoreboard) + 1,
        Username: username,
        Score:    score,
    })
}
