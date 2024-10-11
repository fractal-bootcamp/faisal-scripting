package main

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "regexp"
    "strings"
)

// Helper function to run external shell commands
func runCommand(name string, arg ...string) error {
    cmd := exec.Command(name, arg...)
    cmd.Stdout = os.Stdout
    cmd.Stdin = os.Stdin
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("failed to execute %s: %w", name, err)
    }
    return nil
}

// Helper function to add a script to package.json
func addPackageScript(dir string, script string) error {
    filePath := dir + "/package.json"
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        return fmt.Errorf("package.json does not exist in %s", dir)
    }

    input, err := os.ReadFile(filePath)
    if err != nil {
        return fmt.Errorf("failed to read package.json: %w", err)
    }

    content := strings.Replace(string(input), `"scripts": {`, `"scripts": {`+"\n    "+script+",", 1)
    if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
        return fmt.Errorf("failed to update package.json: %w", err)
    }
    fmt.Println("Successfully updated package.json with the new script.")
    return nil
}

// Helper function to get a yes/no response from the user
func getYesNoResponse(prompt string) bool {
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Print(prompt + " (y/n): ")
        response, _ := reader.ReadString('\n')
        response = strings.TrimSpace(strings.ToLower(response))

        if response == "y" || response == "" { // Treat empty input as "yes"
            return true
        } else if response == "n" {
            return false
        } else {
            fmt.Println("Invalid input. Please enter 'y' for yes or 'n' for no.")
        }
    }
}

// Helper function to create a directory
func createDirectory(dirName string) error {
    if err := os.Mkdir(dirName, 0755); err != nil {
        return fmt.Errorf("failed to create %s directory: %w", dirName, err)
    }
    return nil
}

// Helper function to handle errors
func handleError(action string, err error) {
    if err != nil {
        fmt.Printf("Error while %s: %v\n", action, err)
    }
}

// Basic URL validation
func isValidURL(url string) bool {
    regex := `^(https?|git)://[^\s/$.?#].[^\s]*$`
    match, _ := regexp.MatchString(regex, url)
    return match
}

// welcome message function
func welcomeMessage() {
    fmt.Println("Welcome to Jeez! Press Enter to start setting up your project.")
    bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// Step 1: Prompt for project name
func promptProjectName() (string, error) {
    for {
        fmt.Print("Enter project name: ")
        reader := bufio.NewReader(os.Stdin)
        projectName, _ := reader.ReadString('\n')
        projectName = strings.TrimSpace(projectName)

        if projectName == "" {
            fmt.Println("Jeez!! Project name cannot be empty. Please try again.")
            continue
        }

        if _, err := os.Stat(projectName); !os.IsNotExist(err) {
            fmt.Printf("A directory with the name '%s' already exists. Please choose a different name.\n", projectName)
            continue
        }

        err := createDirectory(projectName)
        if err != nil {
            return "", fmt.Errorf("failed to create project directory: %w", err)
        }

        if err := os.Chdir(projectName); err != nil {
            return "", fmt.Errorf("failed to change to project directory: %w", err)
        }

        return projectName, nil
    }
}

// Step 2: Initialize Git repository
func initializeGit(projectName string) error {
    if !getYesNoResponse("Do you want to initialize a Git repository") {
        fmt.Println("Skipping Git initialization.")
        return nil
    }

    fmt.Println("Initializing git repo...")
    if err := runCommand("git", "init"); err != nil {
        return fmt.Errorf("failed to initialize git repository: %w", err)
    }
    readmeContent := fmt.Sprintf("# %s", projectName)
    if err := os.WriteFile("README.md", []byte(readmeContent), 0644); err != nil {
        return fmt.Errorf("failed to create README.md: %w", err)
    }
    if err := runCommand("git", "add", "README.md"); err != nil {
        return fmt.Errorf("failed to add README.md to git: %w", err)
    }
    if err := runCommand("git", "commit", "-m", "Initial commit for "+projectName); err != nil {
        return fmt.Errorf("failed to commit README.md: %w", err)
    }
    fmt.Println("Jeez! Git repository initialized successfully.")
    return nil
}

// Step 3: Create frontend and backend directories
func createDirectories() (bool, bool) {
    reader := bufio.NewReader(os.Stdin)
    frontend := false
    backend := false

    for {
        fmt.Println("Select your project setup:")
        fmt.Println("1. Frontend")
        fmt.Println("2. Backend")
        fmt.Println("3. Fullstack")
        fmt.Print("Enter your choice (1/2/3): ")
        choice, _ := reader.ReadString('\n')
        choice = strings.TrimSpace(choice)

        switch choice {
        case "1":
            frontend = true
            fmt.Println("Setting up Frontend directory...")
            handleError("creating frontend directory", createDirectory("frontend"))
            return frontend, backend
        case "2":
            backend = true
            fmt.Println("Setting up Backend directory...")
            handleError("creating backend directory", createDirectory("backend"))
            return frontend, backend
        case "3":
            frontend = true
            backend = true
            fmt.Println("Setting up Frontend and Backend directories...")
            handleError("creating frontend directory", createDirectory("frontend"))
            handleError("creating backend directory", createDirectory("backend"))
            return frontend, backend
        default:
            fmt.Println("Invalid choice. Please select 1, 2, or 3.")
        }
    }
}

// Step 4: Set up frontend
func setupFrontend() error {
    reader := bufio.NewReader(os.Stdin)

    for {
        fmt.Println("Select your frontend setup:")
        fmt.Println("1. Vite")
        fmt.Println("2. Skip")
        fmt.Print("Enter your choice (1/2): ")
        choice, _ := reader.ReadString('\n')
        choice = strings.TrimSpace(choice)

        switch choice {
        case "1":
            if err := os.Chdir("frontend"); err != nil {
                return fmt.Errorf("failed to change to frontend directory: %w", err)
            }
            if err := runCommand("npm", "create", "vite@latest", "."); err != nil {
                os.Chdir("..")
                return fmt.Errorf("failed to create Vite project: %w", err)
            }
            if err := runCommand("npm", "install"); err != nil {
                os.Chdir("..")
                return fmt.Errorf("failed to install dependencies: %w", err)
            }
            if err := addPackageScript("frontend", `"dev": "vite"`); err != nil {
                fmt.Println("Warning: Failed to add 'dev' script to package.json:", err)
            }
            fmt.Println("Jeez! Frontend setup complete.")
            os.Chdir("..")
            return nil
        case "2":
            fmt.Println("Skipping frontend setup.")
            return nil
        default:
            fmt.Println("Invalid choice. Please enter 1 or 2.")
        }
    }
}

// Step 5: Set up backend
func setupBackend() error {
    reader := bufio.NewReader(os.Stdin)

    for {
        fmt.Println("Select your backend setup:")
        fmt.Println("1. Express (with TypeScript)")
        fmt.Println("2. Skip")
        fmt.Print("Enter your choice (1/2): ")
        choice, _ := reader.ReadString('\n')
        choice = strings.TrimSpace(choice)

        if choice == "1" {
            if err := os.Chdir("backend"); err != nil {
                handleError("changing to backend directory", err)
                return err
            }
            // Create package.json file if it doesn't exist
            if err := runCommand("npm", "init", "-y"); err != nil {
                os.Chdir("..")
                handleError("initializing npm", err)
                return err
            }
            // Install necessary dependencies
            if err := runCommand("npm", "install", "express", "typescript", "@types/express", "ts-node", "nodemon"); err != nil {
                os.Chdir("..")
                handleError("installing backend dependencies", err)
                return err
            }
            // Create server.ts file
            if err := os.WriteFile("server.ts", []byte(`
import express, { Request, Response } from 'express';

const app = express();
const PORT = 3000;

app.use(express.json());

app.get('/', (req: Request, res: Response) => {
    res.send('Hello, Jeez!');
});

app.listen(PORT, () => {
    console.log("Server is running on http://localhost:" + PORT);
});
`), 0644); err != nil {
                os.Chdir("..")
                handleError("creating server.ts file", err)
                return err
            }
            // Add the start script to package.json
            if err := addPackageScript("backend", `"start": "nodemon server.ts"`); err != nil {
                fmt.Println("Warning: Failed to add 'start' script to package.json:", err)
            }
            fmt.Println("Jeez! Backend setup complete.")
            os.Chdir("..")
            return nil
        } else if choice == "2" {
            fmt.Println("Skipping backend setup.")
            return nil
        } else {
            fmt.Println("Invalid choice. Please enter 1 or 2.")
        }
    }
}

// Step 6: Add Docker and PostgreSQL setup
func setupDatabase(dirName string) error {
    if err := os.Chdir("backend"); err != nil {
        return fmt.Errorf("failed to change to backend directory: %w", err)
    }

    if !getYesNoResponse("Do you want to set up a database") {
        fmt.Println("Skipping database setup.")
        return nil
    }

    // Use the dirName to set the POSTGRES_DB name
    dockerComposeContent := fmt.Sprintf(`version: '3.8'
services:
    postgres:
      image: postgres:13
      environment:
        POSTGRES_USER: postgres
        POSTGRES_PASSWORD: postgres
        POSTGRES_DB: %s_db
      command: -c fsync=off -c full_page_writes=off -c synchronous_commit=off -c max_connections=500
      ports:
        - 10001:5432
    `, dirName)

    if err := os.WriteFile("docker-compose.yml", []byte(dockerComposeContent), 0644); err != nil {
        os.Chdir("..")
        return fmt.Errorf("failed to write docker-compose.yml: %w", err)
    }
    fmt.Println("Jeez! Database setup complete.")
    os.Chdir("..")
    return nil
}

// Step 7: Set up ORM for backend
func setupOrm() error {
    if !getYesNoResponse("Do you want to set up an ORM (Prisma)") {
        fmt.Println("Skipping ORM setup.")
        return nil
    }

    if _, err := os.Stat("backend"); os.IsNotExist(err) {
        return fmt.Errorf("backend directory does not exist. Skipping ORM setup.")
    }
    if err := os.Chdir("backend"); err != nil {
        return fmt.Errorf("failed to change to backend directory: %w", err)
    }

    setupSteps := []struct {
        command  string
        args     []string
        errorMsg string
    }{
        {"npx", []string{"prisma", "init"}, "Failed to initialize Prisma"},
        {"npm", []string{"install", "@prisma/client"}, "Failed to install Prisma client"},
    }

    for _, step := range setupSteps {
        if err := runCommand(step.command, step.args...); err != nil {
            os.Chdir("..")
            return fmt.Errorf("%s: %w", step.errorMsg, err)
        }
    }

    // Add datasource and basic model to schema.prisma
    prismaSchemaPath := "prisma/schema.prisma"
    datasourceAndModel := `
datasource db {
    provider = "postgresql"
    url      = env("DATABASE_URL")
}

model User {
    id        String  @id @default(uuid())
    email     String  @unique
    firstName String
    lastName  String
}
    `
    if err := os.WriteFile(prismaSchemaPath, []byte(datasourceAndModel), 0644); err != nil {
        os.Chdir("..")
        return fmt.Errorf("failed to write to schema.prisma: %w", err)
    }

    // Run Prisma generate
    if err := runCommand("npx", "prisma", "generate"); err != nil {
        os.Chdir("..")
        return fmt.Errorf("failed to generate Prisma client: %w", err)
    }

    fmt.Println("Jeez! ORM setup complete.")
    os.Chdir("..")
    return nil
}

// Step 8: Update .env.local file
func updateEnvFile() error {
    if !getYesNoResponse("Do you want to create an .env.local file for database configuration") {
        fmt.Println("Skipping .env.local setup.")
        return nil
    }

    envContent := "DATABASE_URL=postgresql://myuser:mypassword@localhost:1001/mydb"
    if err := os.WriteFile(".env.local", []byte(envContent), 0644); err != nil {
        return fmt.Errorf("failed to create .env.local file: %w", err)
    }
    fmt.Println("Jeez! .env.local file created.")
    return nil
}

// Step 9: Set up remote Git repository
func setupGitRemote() error {
    if !getYesNoResponse("Do you want to set up a remote Git repository") {
        fmt.Println("Skipping remote Git setup.")
        return nil
    }

    if _, err := os.Stat(".git"); os.IsNotExist(err) {
        return fmt.Errorf("Git repository is not initialized. Please initialize Git before setting up a remote.")
    }

    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("Enter the remote repo URL: ")
        remoteUrl, _ := reader.ReadString('\n')
        remoteUrl = strings.TrimSpace(remoteUrl)

        if !isValidURL(remoteUrl) {
            fmt.Println("Invalid URL format. Please enter a valid URL.")
            continue
        }

        if err := runCommand("git", "remote", "add", "origin", remoteUrl); err != nil {
            fmt.Println("Failed to add remote origin:", err)
            fmt.Println("Please re-enter the remote repo URL or type 'skip' to skip this step.")
            if remoteUrl == "skip" {
                return nil
            }
            continue
        }

        if err := runCommand("git", "branch", "-M", "main"); err != nil {
            return fmt.Errorf("failed to rename branch to main: %w", err)
        }
        if err := runCommand("git", "push", "-u", "origin", "main"); err != nil {
            return fmt.Errorf("failed to push to remote repository: %w", err)
        }
        fmt.Println("Jeez! Git remote repo complete.")
        break
    }
    return nil
}

////// Main script //////
func main() {
    welcomeMessage()
    projectName, err := promptProjectName()
    if err != nil {
        handleError("prompting project name", err)
        return
    }

    initializeGit(projectName)
    frontend, backend := createDirectories()

    if frontend {
        if err := setupFrontend(); err != nil {
            handleError("setting up frontend", err)
        }
    }
    if backend {
        if err := setupBackend(); err != nil {
            handleError("setting up backend", err)
        }
        if err := setupDatabase(projectName); err != nil { // Pass projectName to setupDatabase
            handleError("setting up database", err)
        }
        if err := setupOrm(); err != nil {
            handleError("setting up ORM", err)
        }
        if err := updateEnvFile(); err != nil {
            handleError("updating .env file", err)
        }
    }
    if err := setupGitRemote(); err != nil {
        handleError("setting up Git remote", err)
    }

    fmt.Println("Jeeeez! Project setup is ready, let's rip some code!")
}