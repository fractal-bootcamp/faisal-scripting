package main

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "regexp"
    "strings"
)

const (
    ColorReset  = "\033[0m"
    ColorRed    = "\033[31m"
    ColorGreen  = "\033[32m"
    ColorYellow = "\033[33m"
    ColorBlue   = "\033[34m"
    ColorPurple = "\033[35m"
    ColorCyan   = "\033[36m"
    ColorBold   = "\033[1m"
)

// Helper function to run external shell commands
func runCommand(name string, arg ...string) error {
    cmd := exec.Command(name, arg...)
    cmd.Stdout = os.Stdout
    cmd.Stdin = os.Stdin
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("%sfailed to execute %s: %w%s", ColorRed, name, err, ColorReset)
    }
    return nil
}

// Helper function to add a script to package.json
func addPackageScript(dir string, script string) error {
    filePath := dir + "/package.json"
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        return fmt.Errorf("%spackage.json does not exist in %s%s", ColorRed, dir, ColorReset)
    }

    input, err := os.ReadFile(filePath)
    if err != nil {
        return fmt.Errorf("%sfailed to read package.json: %w%s", ColorRed, err, ColorReset)
    }

    content := strings.Replace(string(input), `"scripts": {`, `"scripts": {`+"\n    "+script+",", 1)
    if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
        return fmt.Errorf("%sfailed to update package.json: %w%s", ColorRed, err, ColorReset)
    }
    fmt.Printf("%sSuccessfully updated package.json with the new script.%s\n", ColorGreen, ColorReset)
    return nil
}

// Helper function to get a yes/no response from the user
func getYesNoResponse(prompt string) bool {
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Printf("%s%s (y/n): %s", ColorYellow, prompt, ColorReset)
        response, _ := reader.ReadString('\n')
        response = strings.TrimSpace(strings.ToLower(response))

        if response == "y" || response == "" { // Treat empty input as "yes"
            return true
        } else if response == "n" {
            return false
        } else {
            fmt.Printf("%sInvalid input. Please enter 'y' for yes or 'n'.%s\n", ColorRed, ColorReset)
        }
    }
}

// Helper function to create a directory
func createDirectory(dirName string) error {
    if err := os.Mkdir(dirName, 0755); err != nil {
        return fmt.Errorf("%sfailed to create %s directory: %w%s", ColorRed, dirName, err, ColorReset)
    }
    return nil
}

// Helper function to handle errors
func handleError(action string, err error) {
    if err != nil {
        fmt.Printf("%sError while %s: %v%s\n", ColorRed, action, err, ColorReset)
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
    fmt.Printf("%s%sWelcome to Jeez! Let's start setting up your project.%s%s\n", ColorBold, ColorCyan, ColorReset, ColorReset)
}

// Step 1: Prompt for project name
func promptProjectName() (string, error) {
    for {
        fmt.Printf("%sEnter project name: %s", ColorYellow, ColorReset)
        reader := bufio.NewReader(os.Stdin)
        projectName, _ := reader.ReadString('\n')
        projectName = strings.TrimSpace(projectName)

        if projectName == "" {
            fmt.Printf("%sJeez!! Project name cannot be empty. Please try again.%s\n", ColorRed, ColorReset)
            continue
        }

        if _, err := os.Stat(projectName); !os.IsNotExist(err) {
            fmt.Printf("%sA directory with the name '%s' already exists. Please choose a different name.%s\n", ColorRed, projectName, ColorReset)
            continue
        }

        err := createDirectory(projectName)
        if err != nil {
            return "", fmt.Errorf("%sfailed to create project directory: %w%s", ColorRed, err, ColorReset)
        }

        if err := os.Chdir(projectName); err != nil {
            return "", fmt.Errorf("%sfailed to change to project directory: %w%s", ColorRed, err, ColorReset)
        }

        fmt.Printf("%sProject '%s' created successfully!%s\n", ColorGreen, projectName, ColorReset)
        return projectName, nil
    }
}

// Step 1.1: Prompt to add Bun
func promptAddBun() error {
    if getYesNoResponse("Do you want to add Bun to your project") {
        if err := runCommand("bun", "init"); err != nil {
            return fmt.Errorf("%sfailed to initialize Bun: %w%s", ColorRed, err, ColorReset)
        }
        fmt.Printf("%sJeez! Bun added to the project successfully.%s\n", ColorGreen, ColorReset)
    } else {
        fmt.Printf("%sSkipping Bun setup.%s\n", ColorYellow, ColorReset)
    }
    return nil
}

// Step 2: Initialize Git repository
func initializeGit(projectName string) error {
    if !getYesNoResponse("Do you want to initialize a Git repository") {
        fmt.Println("Skipping Git initialization.") //////////// add color
        return nil
    }

    fmt.Println("Initializing git repo...")
    if err := runCommand("git", "init"); err != nil {
        return fmt.Errorf("%sfailed to initialize git repository: %w%s", ColorRed, err, ColorReset)
    }
    readmeContent := fmt.Sprintf("# %s", projectName)
    if err := os.WriteFile("README.md", []byte(readmeContent), 0644); err != nil {
        return fmt.Errorf("%sfailed to create README.md: %w%s", ColorRed, err, ColorReset)
    }
    if err := runCommand("git", "add", "README.md"); err != nil {
        return fmt.Errorf("%sfailed to add README.md to git: %w%s", ColorRed, err, ColorReset)
    }
    if err := runCommand("git", "commit", "-m", "Initial commit for "+projectName); err != nil {
        return fmt.Errorf("%sfailed to commit README.md: %w%s", ColorRed, err, ColorReset)
    }
    fmt.Printf("%sJeez! Git repository initialized successfully.%s\n", ColorGreen, ColorReset)
    return nil
}

// Step 3: Create frontend and backend directories
func createDirectories() (bool, bool) {
    reader := bufio.NewReader(os.Stdin)
    frontend := false
    backend := false

    for {
        fmt.Printf("%sSelect your project setup:%s\n", ColorBold, ColorReset)
        fmt.Println("1. Frontend")
        fmt.Println("2. Backend")
        fmt.Println("3. Fullstack")
        fmt.Printf("%sEnter your choice (1/2/3): %s", ColorYellow, ColorReset)
        choice, _ := reader.ReadString('\n')
        choice = strings.TrimSpace(choice)

        switch choice {
        case "1":
            frontend = true
            fmt.Printf("%sSetting up Frontend directory...%s\n", ColorBlue, ColorReset)
            handleError("creating frontend directory", createDirectory("frontend"))
            return frontend, backend
        case "2":
            backend = true
            fmt.Printf("%sSetting up Backend directory...%s\n", ColorBlue, ColorReset)
            handleError("creating backend directory", createDirectory("backend"))
            return frontend, backend
        case "3":
            frontend = true
            backend = true
            fmt.Printf("%sSetting up Frontend and Backend directories...%s\n", ColorBlue, ColorReset)
            handleError("creating frontend directory", createDirectory("frontend"))
            handleError("creating backend directory", createDirectory("backend"))
            return frontend, backend
        default:
            fmt.Printf("%sInvalid choice. Please select 1, 2, or 3.%s\n", ColorRed, ColorReset)
        }
    }
}

// Step 4: Set up frontend
func setupFrontend() error {
    reader := bufio.NewReader(os.Stdin)

    for {
        fmt.Printf("%sSelect your frontend setup:%s\n", ColorBold, ColorReset)
        fmt.Println("1. Vite")
        fmt.Println("2. Skip")
        fmt.Printf("%sEnter your choice (1/2): %s", ColorYellow, ColorReset)
        choice, _ := reader.ReadString('\n')
        choice = strings.TrimSpace(choice)

        switch choice {
        case "1":
            if err := os.Chdir("frontend"); err != nil {
                return fmt.Errorf("%sfailed to change to frontend directory: %w%s", ColorRed, err, ColorReset)
            }
            if err := runCommand("npm", "create", "vite@latest", "."); err != nil {
                os.Chdir("..")
                return fmt.Errorf("%sfailed to create Vite project: %w%s", ColorRed, err, ColorReset)
            }
            if err := runCommand("npm", "install"); err != nil {
                os.Chdir("..")
                return fmt.Errorf("%sfailed to install dependencies: %w%s", ColorRed, err, ColorReset)
            }
            if err := addPackageScript("frontend", `"dev": "vite"`); err != nil {
                fmt.Printf("%sWarning: Failed to add 'dev' script to package.json:%s %v\n", ColorYellow, ColorReset, err)
            }
            fmt.Printf("%sJeez! Vite setup complete.%s\n", ColorGreen, ColorReset)

            // Ask if user wants to install TailwindCSS
            if getYesNoResponse("Do you want to install TailwindCSS for Vite-React") {
                fmt.Printf("%sInstalling TailwindCSS...%s\n", ColorBlue, ColorReset)
                if err := runCommand("npm", "install", "-D", "tailwindcss", "postcss", "autoprefixer"); err != nil {
                    fmt.Printf("%sFailed to install TailwindCSS: %v%s\n", ColorRed, err, ColorReset)
                } else {
                    if err := runCommand("npx", "tailwindcss", "init", "-p"); err != nil {
                        fmt.Printf("%sFailed to initialize TailwindCSS: %v%s\n", ColorRed, err, ColorReset)
                    } else {
                        fmt.Printf("%sTailwindCSS setup complete.%s\n", ColorGreen, ColorReset)
                    }
                }
            }

            // Ask if user wants to install Storybook
            if getYesNoResponse("Do you want to install Storybook") {
                fmt.Printf("%sInstalling Storybook...%s\n", ColorBlue, ColorReset)
                if err := runCommand("npx", "storybook", "init"); err != nil {
                    fmt.Printf("%sFailed to install Storybook: %v%s\n", ColorRed, err, ColorReset)
                } else {
                    fmt.Printf("%sStorybook setup complete.%s\n", ColorGreen, ColorReset)
                }
            }

            os.Chdir("..")
            return nil
        case "2":
            fmt.Printf("%sSkipping frontend setup.%s\n", ColorYellow, ColorReset)
            return nil
        default:
            fmt.Printf("%sInvalid choice. Please enter 1 or 2.%s\n", ColorRed, ColorReset)
        }
    }
}

// Step 5: Set up backend
func setupBackend() error {
    reader := bufio.NewReader(os.Stdin)

    for {
        fmt.Printf("%sSelect your backend setup:%s\n", ColorBold, ColorReset)
        fmt.Println("1. Express (with TypeScript)")
        fmt.Println("2. Skip")
        fmt.Printf("%sEnter your choice (1/2): %s", ColorYellow, ColorReset)
        choice, _ := reader.ReadString('\n')
        choice = strings.TrimSpace(choice)

        if choice == "1" {
            if err := os.Chdir("backend"); err != nil {
                handleError("changing to backend directory", err)
                return err
            }

            // Initialize npm and create package.json if it doesn't exist
            if err := runCommand("npm", "init", "-y"); err != nil {
                os.Chdir("..")
                handleError("initializing npm", err)
                return err
            }

            // Install necessary dependencies including dotenv and CORS
            if err := runCommand("npm", "install", "express", "typescript", "@types/express", "ts-node", "nodemon", "cors", "@types/cors", "dotenv"); err != nil {
                os.Chdir("..")
                handleError("installing backend dependencies", err)
                return err
            }

            // Create tsconfig.json file for TypeScript configuration
            tsConfigContent := `{
"compilerOptions": {
    "target": "ES6",
    "module": "commonjs",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "outDir": "./dist",
    "rootDir": "./src"
    },
    "include": ["src/**/*.ts"],
    "exclude": ["node_modules"]
}`
            if err := os.WriteFile("tsconfig.json", []byte(tsConfigContent), 0644); err != nil {
                os.Chdir("..")
                handleError("creating tsconfig.json file", err)
                return err
            }

            // Create a simple Express server with CORS and dotenv support in TypeScript
            if err := os.Mkdir("src", 0755); err != nil {
                os.Chdir("..")
                handleError("creating src directory", err)
                return err
            }

            serverContent := `import express, { Request, Response } from 'express';
import { PrismaClient } from "@prisma/client";
import cors from 'cors';
import dotenv from 'dotenv';

dotenv.config();

const app = express();
const PORT = process.env.PORT || 3000;

app.use(cors({ origin: "*" }));
app.use(express.json());

app.get('/', (req: Request, res: Response) => {
    res.send('Hello, Jeez!');
});

app.listen(PORT, () => {
    console.log("Server is running on http://localhost:" + PORT);
});
`
            if err := os.WriteFile("src/server.ts", []byte(serverContent), 0644); err != nil {
                os.Chdir("..")
                handleError("creating server.ts file", err)
                return err
            }

            // Add start and build scripts to package.json
            if err := addPackageScript("backend", `"start": "nodemon src/server.ts", "build": "tsc"`); err != nil {
                fmt.Printf("%sWarning: Failed to add 'start' and 'build' scripts to package.json:%s %v\n", ColorYellow, ColorReset, err)
            }

            fmt.Printf("%sJeez! Backend setup with TypeScript, CORS, and dotenv complete.%s\n", ColorGreen, ColorReset)
            os.Chdir("..")
            return nil
        } else if choice == "2" {
            fmt.Println("Skipping backend setup.") //////////// add color
            return nil
        } else {
            fmt.Printf("%sInvalid choice. Please enter 1 or 2.%s\n", ColorRed, ColorReset)
        }
    }
}

// Step 6: Add Docker and PostgreSQL setup
func setupDatabase(dirName string) error {
    if err := os.Chdir("backend"); err != nil {
        return fmt.Errorf("%sfailed to change to backend directory: %w%s", ColorRed, err, ColorReset)
    }

    if !getYesNoResponse("Do you want to set up a database") {
        fmt.Println("Skipping database setup.") //////////// add color
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
        return fmt.Errorf("%sfailed to write docker-compose.yml: %w%s", ColorRed, err, ColorReset)
    }
    fmt.Printf("%sJeez! Database setup complete.%s\n", ColorGreen, ColorReset)
    os.Chdir("..")
    return nil
}

// Step 7: Set up ORM for backend
func setupOrm() error {
    if !getYesNoResponse("Do you want to set up an ORM (Prisma)") {
        fmt.Printf("%sSkipping ORM setup.%s\n", ColorYellow, ColorReset)
        return nil
    }

    if _, err := os.Stat("backend"); os.IsNotExist(err) {
        return fmt.Errorf("%sBackend directory does not exist. Skipping ORM setup.%s", ColorRed, ColorReset)
    }
    if err := os.Chdir("backend"); err != nil {
        return fmt.Errorf("%sFailed to change to backend directory: %w%s", ColorRed, err, ColorReset)
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
            return fmt.Errorf("%s%s: %w%s", ColorRed, step.errorMsg, err, ColorReset)
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
        return fmt.Errorf("%sFailed to write to schema.prisma: %w%s", ColorRed, err, ColorReset)
    }

    // Run Prisma generate
    if err := runCommand("npx", "prisma", "generate"); err != nil {
        os.Chdir("..")
        return fmt.Errorf("%sFailed to generate Prisma client: %w%s", ColorRed, err, ColorReset)
    }

    fmt.Printf("%sJeez! ORM setup complete.%s\n", ColorGreen, ColorReset)
    os.Chdir("..")
    return nil
}

// Step 8: Update .env.local file
func updateEnvFile(dirName string) error {
    if !getYesNoResponse("Do you want to create an .env.local file for database configuration") {
        fmt.Println("Skipping .env.local setup.") //////////// add color
        return nil
    }

    // Update DATABASE_URL to reflect the PostgreSQL setup in Docker
    envContent := fmt.Sprintf("DATABASE_URL=postgresql://postgres:postgres@localhost:10001/%s_db\nPORT=3000", dirName)
    if err := os.WriteFile("backend/.env.local", []byte(envContent), 0644); err != nil {
        return fmt.Errorf("%sfailed to create .env.local file: %w%s", ColorRed, err, ColorReset)
    }
    fmt.Printf("%sJeez! .env.local file created with database configuration.%s\n", ColorGreen, ColorReset)
    return nil
}

// Step 9: Set up remote Git repository
func setupGitRemote() error {
    if !getYesNoResponse("Do you want to set up a remote Git repository") {
        fmt.Println("Skipping remote Git setup.")
        return nil
    }

    if _, err := os.Stat(".git"); os.IsNotExist(err) {
        return fmt.Errorf("%sGit repository is not initialized. Please initialize Git before setting up a remote.%s", ColorRed, ColorReset)
    }

    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Printf("%sEnter the remote repo URL: %s", ColorYellow, ColorReset)
        remoteUrl, _ := reader.ReadString('\n')
        remoteUrl = strings.TrimSpace(remoteUrl)

        if !isValidURL(remoteUrl) {
            fmt.Printf("%sInvalid URL format. Please enter a valid URL.%s\n", ColorRed, ColorReset)
            continue
        }

        if err := runCommand("git", "remote", "add", "origin", remoteUrl); err != nil {
            fmt.Printf("%sFailed to add remote origin:%s %v\n", ColorRed, ColorReset, err)
            fmt.Println("Please re-enter the remote repo URL or type 'skip' to skip this step.")
            if remoteUrl == "skip" {
                return nil
            }
            continue
        }

        if err := runCommand("git", "branch", "-M", "main"); err != nil {
            return fmt.Errorf("%sfailed to rename branch to main: %w%s", ColorRed, err, ColorReset)
        }
        if err := runCommand("git", "push", "-u", "origin", "main"); err != nil {
            return fmt.Errorf("%sfailed to push to remote repository: %w%s", ColorRed, err, ColorReset)
        }
        fmt.Printf("%sJeez! Git remote repo complete.%s\n", ColorGreen, ColorReset)
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

    // Prompt to add Bun
    if err := promptAddBun(); err != nil {
        handleError("adding Bun", err)
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
        if err := updateEnvFile(projectName); err != nil {
            handleError("updating .env file", err)
        }
    }
    if err := setupGitRemote(); err != nil {
        handleError("setting up Git remote", err)
    }

    fmt.Printf("%s%sJeeeez! Project setup is ready, let's rip some code!%s%s\n", ColorBold, ColorGreen, ColorReset, ColorReset)
}