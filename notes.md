// context/AuthContext.jsx  !!adding token generation token assignment to the frontend
import { createContext, useState, useEffect } from 'react';
import { api } from '../api/client'; // your api client

export const AuthContext = createContext();

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [guestId, setGuestId] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    restoreSession();
  }, []);

  // 1. Get or create guest token
  const getGuestToken = () => {
    const saved = localStorage.getItem('guestId');
    if (saved) return saved;
    const id = crypto.randomUUID();
    localStorage.setItem('guestId', id);
    return id;
  };

  // 2. On load: restore user from backend
  const restoreSession = async () => {
    const token = getGuestToken();
    setGuestId(token);

    try {
      // Call backend: "Who is this guest_token?"
      const response = await api.get('/api/session');
      // Response: { user: { id, email }, is_guest: false }
      setUser(response.user);
    } catch (err) {
      // Maybe no user yet — that’s fine
      console.log("No registered user found");
    } finally {
      setLoading(false);
    }
  };

  // 3. Login: user registers
  const login = async (email, password) => {
    const response = await api.post('/api/register', { email, password });
    setUser(response.user); // now they're real
    // token already in localStorage
  };

  const logout = () => {
    localStorage.removeItem('token');
    setUser(null);
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <AuthContext.Provider value={{
      user,
      guestId,
      login,
      logout,
    }}>
      {children}
    </AuthContext.Provider>
  );
}




// handlers/session.go
func SessionHandler(pool *pgxpool.Pool) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        guestToken := r.Header.Get("X-Guest-Token")
        if guestToken == "" {
            http.Error(w, "Missing guest token", http.StatusBadRequest)
            return
        }

        var (
            userID   string
            email    *string
            isGuest  bool
        )

        err := pool.QueryRow(r.Context(), `
            SELECT id, email, is_guest
            FROM users
            WHERE guest_token = $1
        `, guestToken).Scan(&userID, &email, &isGuest)

        if err != nil {
            if err == pgx.ErrNoRows {
                // Auto-create guest
                newID := "usr_" + utils.RandomString(16)
                _, err := pool.Exec(r.Context(), `
                    INSERT INTO users (id, guest_token, is_guest, created_at)
                    VALUES ($1, $2, true, NOW())
                `, newID, guestToken)
                if err != nil {
                    http.Error(w, "Failed to create guest", http.StatusInternalServerError)
                    return
                }
                userID = newID
                isGuest = true
            } else {
                http.Error(w, "Server error", http.StatusInternalServerError)
                return
            }
        }

        // Respond
        var userResponse map[string]interface{}
        if email != nil {
            userResponse = map[string]interface{}{"id": userID, "email": *email}
        }

        respondJSON(w, map[string]interface{}{
            "user":     userResponse,
            "is_guest": isGuest,
            "guest_token": guestToken,
        })
    }
}



// api/client.js
const API_BASE = import.meta.env.VITE_API_URL;

export const api = {
  get: async (url) => {
    const guestToken = localStorage.getItem('guestId') || crypto.randomUUID();
    const response = await fetch(`${API_BASE}${url}`, {
      headers: {
        'Content-Type': 'application/json',
        'X-Guest-Token': guestToken,
      },
    });
    return response.json();
  },
  post: async (url, data) => {
    const guestToken = localStorage.getItem('guestId') || crypto.randomUUID();
    const response = await fetch(`${API_BASE}${url}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Guest-Token': guestToken,
      },
      body: JSON.stringify(data),
    });
    return response.json();
  },
};



// internal/data/load_combos.go
package data

import (
    "bufio"
    "context"
    "fmt"
    "os"
    "strings"

    "github.com/jackc/pgx/v5/pgxpool"
)

// LoadCombosFromFile reads cleaned combos from file and inserts into DB !! load combos from file load combination from file
func LoadCombosFromFile(pool *pgxpool.Pool, filePath string) error {
    file, err := os.Open(filePath)
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    // Use transaction for atomicity
    tx, err := pool.Begin(context.Background())
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback(context.Background())

    scanner := bufio.NewScanner(file)
    count := 0

    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" {
            continue
        }

        // Split into parts (handles multiple spaces)
        parts := strings.Fields(line)
        if len(parts) != 4 {
            return fmt.Errorf("invalid combo (expected 4 numbers): %s", line)
        }

        // Create ID: "8-9-11-12"
        id := strings.Join(parts, "-")

        // Insert into DB
        _, err := tx.Exec(context.Background(),
            `INSERT INTO combos (id, cards) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING`,
            id, parts, // parts is []string → becomes PostgreSQL TEXT[]
        )
        if err != nil {
            return fmt.Errorf("failed to insert combo %s: %w", id, err)
        }

        count++
    }

    if err := scanner.Err(); err != nil {
        return fmt.Errorf("error reading file: %w", err)
    }

    // Commit transaction
    if err := tx.Commit(context.Background()); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    fmt.Printf("✅ Loaded %d combos into the database\n", count)
    return nil
}

// internal/migrations/migrate.go   how to use the above code
package migrations

import (
    "context"
    "embed"

    "github.com/pressly/goose/v3"
    "github.com/yourname/sanbar/db"
    "github.com/yourname/sanbar/data"
)

//go:embed *.sql
var embedMigrations embed.FS

func Migrate(pool *pgxpool.Pool) error {
    goose.SetBaseFS(embedMigrations)
    if err := goose.SetDialect("postgres"); err != nil {
        return err
    }

    // Run schema migrations
    if err := goose.Up(context.Background(), &db.MigrationConn{Pool: pool}, "."); err != nil {
        return fmt.Errorf("goose up: %w", err)
    }

    // Load combos data
    if err := data.LoadCombosFromFile(pool, "./data/clean.txt"); err != nil {
        return fmt.Errorf("failed to load combos: %w", err)
    }

    return nil
}
