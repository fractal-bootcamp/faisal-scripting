import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

// welcome message function at the start of the script
func welcomeMessage() {
	fmt.Println("Welcome to Jeez! Press Enter to start setting up your project.")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// 1. Prompt for project name
func promptProjectName() string {
    for {
        fmt.Print("Enter project name: ")
        reader := bufio.NewReader(os.Stdin)
        projectName, _ := reader.ReadString('\n')
        projectName = strings.TrimSpace(projectName)

        if projectName == "" {
            fmt.Println("Jeez!! Project name cannot be empty. Please try again.")
            continue // Re-prompt the user for the project name
        }

        if _, err := os.Stat(projectName); !os.IsNotExist(err) {
            fmt.Printf("A directory with the name '%s' already exists. Please choose a different name.\n", projectName)
            continue // Re-prompt the user
        }

        err := os.Mkdir(projectName, 0755)
        if err != nil {
            fmt.Println("Yikes! Failed to create project directory:", err)
            os.Exit(1)
        }

        if err := os.Chdir(projectName); err != nil {
            fmt.Println("Failed to change to project directory:", err)
            os.Exit(1)
        }
        return projectName
    }
}

// 2. Initialize Git repository
func initializeGit(projectName string) {
    fmt.Println("Initializing git repo...")
    if err := runCommand("git", "init"); err != nil {
        fmt.Println("Failed to initialize git repository:", err)
        return
    }
    readmeContent := fmt.Sprintf("# %s", projectName)
    if err := os.WriteFile("README.md", []byte(readmeContent), 0644); err != nil {
        fmt.Println("Failed to create README.md:", err)
        return
    }
    if err := runCommand("git", "add", "README.md"); err != nil {
        fmt.Println("Failed to add README.md to git:", err)
        return
    }
    if err := runCommand("git", "commit", "-m", "Initial commit for "+projectName); err != nil {
        fmt.Println("Failed to commit README.md:", err)
        return
    }
    fmt.Println("Jeez! Git repository initialized successfully.")
}

// 3. Create frontend and backend directories
func createDirectories() {
	fmt.Println("Setting up project directories...")
	if err := os.Mkdir("frontend", 0755); err != nil {
		fmt.Println("Failed to create frontend directory:", err)
	}
	if err := os.Mkdir("backend", 0755); err != nil {
		fmt.Println("Failed to create backend directory:", err)
	}
	fmt.Println("Jeez! Directories created.")
}

// 4. Set up frontend
func setupFrontend() error {
    fmt.Print("Choose your frontend setup (Vite-React or skip): ")
    reader := bufio.NewReader(os.Stdin)
    setupFrontend, _ := reader.ReadString('\n')
    setupFrontend = strings.TrimSpace(setupFrontend)

    if setupFrontend == "Vite-React" {
        if err := os.Chdir("frontend"); err != nil {
            return fmt.Errorf("failed to change to frontend directory: %w", err)
        }
        if err := runCommand("npm", "create", "vite@latest", ".", "--", "--template", "react"); err != nil {
            os.Chdir("..")
            return fmt.Errorf("failed to create Vite-React project: %w", err)
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
    } else {
        fmt.Println("Skipping frontend setup.")
    }
    return nil
}

// 5. Set up backend
func setupBackend() error {
    fmt.Print("Choose your backend setup (Express or skip): ")
    reader := bufio.NewReader(os.Stdin)
    setupBackend, _ := reader.ReadString('\n')
    setupBackend = strings.TrimSpace(setupBackend)

    if setupBackend == "Express" {
        if err := os.Chdir("backend"); err != nil {
            return fmt.Errorf("failed to change to backend directory: %w", err)
        }
        if err := runCommand("npm", "init", "-y"); err != nil {
            os.Chdir("..")
            return fmt.Errorf("failed to initialize npm: %w", err)
        }
        if err := runCommand("npm", "install", "express"); err != nil {
            os.Chdir("..")
            return fmt.Errorf("failed to install Express: %w", err)
        }
        if err := addPackageScript("backend", `"start": "node index.js"`); err != nil {
            fmt.Println("Warning: Failed to add 'start' script to package.json:", err)
        }
        fmt.Println("Jeez! Backend setup complete.")
        os.Chdir("..")
    } else {
        fmt.Println("Skipping backend setup.")
    }
    return nil
}

// 6. Add Docker and PostgreSQL setup
func setupDatabase() error {
    fmt.Print("Do you want to set up a database? (y/n): ")
    reader := bufio.NewReader(os.Stdin)
    setupDb, _ := reader.ReadString('\n')
    setupDb = strings.TrimSpace(setupDb)

    if setupDb == "y" {
        dockerComposeContent := `version: '3.8'
services:
  db:
    image: postgres
    environment:
      POSTGRES_USER: myuser
      POSTGRES_PASSWORD: mypassword
      POSTGRES_DB: mydb
    ports:
      - "1001:5432"
`
        if err := os.WriteFile("docker-compose.yml", []byte(dockerComposeContent), 0644); err != nil {
            return fmt.Errorf("failed to write docker-compose.yml: %w", err)
        }
        fmt.Println("Jeez! Database setup complete.")
    } else {
        fmt.Println("Skipping database setup.")
    }
    return nil
}

// 7. Set up ORM for backend
func setupOrm() {
    fmt.Print("Do you want to set up an ORM (Prisma)? (y/n): ")
    reader := bufio.NewReader(os.Stdin)
    setupOrm, _ := reader.ReadString('\n')
    setupOrm = strings.TrimSpace(setupOrm)

    if setupOrm == "y" {
        if _, err := os.Stat("backend"); os.IsNotExist(err) {
            fmt.Println("Backend directory does not exist. Skipping ORM setup.")
            return
        }
        if err := os.Chdir("backend"); err != nil {
            fmt.Println("Failed to change to backend directory:", err)
            return
        }

        setupSteps := []struct {
            command string
            args    []string
            errorMsg string
        }{
            {"npx", []string{"prisma", "init"}, "Failed to initialize Prisma"},
            {"npm", []string{"install", "@prisma/client"}, "Failed to install Prisma client"},
            {"npx", []string{"prisma", "generate"}, "Failed to generate Prisma client"},
        }

        for _, step := range setupSteps {
            if err := runCommand(step.command, step.args...); err != nil {
                fmt.Println(step.errorMsg, ":", err)
                os.Chdir("..")
                return
            }
        }

        fmt.Println("Jeez! ORM setup complete.")
        os.Chdir("..")
    } else {
        fmt.Println("Skipping ORM setup.")
    }
}

// 8. Update .env.local file
func updateEnvFile() error {
    fmt.Print("Do you want to create an .env.local file for database configuration? (y/n): ")
    reader := bufio.NewReader(os.Stdin)
    setupEnv, _ := reader.ReadString('\n')
    setupEnv = strings.TrimSpace(setupEnv)

    if setupEnv == "y" {
        envContent := "DATABASE_URL=postgresql://myuser:mypassword@localhost:1001/mydb"
        if err := os.WriteFile(".env.local", []byte(envContent), 0644); err != nil {
            return fmt.Errorf("failed to create .env.local file: %w", err)
        }
        fmt.Println("Jeez! .env.local file created.")
    } else {
        fmt.Println("Skipping .env.local setup.")
    }
    return nil
}

// 9. push to remote Git repository
func setupGitRemote() error {
    fmt.Print("Do you want to set up a remote Git repository? (y/n): ")
    reader := bufio.NewReader(os.Stdin)
    setupGit, _ := reader.ReadString('\n')
    setupGit = strings.TrimSpace(setupGit)

    if setupGit == "y" {
        if _, err := os.Stat(".git"); os.IsNotExist(err) {
            return fmt.Errorf("git repository is not initialized. Please initialize Git before setting up a remote")
        }

        fmt.Print("Enter the remote repo URL: ")
        remoteUrl, _ := reader.ReadString('\n')
        remoteUrl = strings.TrimSpace(remoteUrl)

        if err := runCommand("git", "remote", "add", "origin", remoteUrl); err != nil {
            return fmt.Errorf("failed to add remote origin: %w", err)
        }
        if err := runCommand("git", "branch", "-M", "main"); err != nil {
            return fmt.Errorf("failed to rename branch to main: %w", err)
        }
        if err := runCommand("git", "push", "-u", "origin", "main"); err != nil {
            return fmt.Errorf("failed to push to remote repository: %w", err)
        }
        fmt.Println("Jeez! Git remote repo complete.")
    } else {
        fmt.Println("Skipping remote Git setup.")
    }
    return nil
}

// Execute external shell commands
func runCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute %s: %w", name, err)
	}
	return nil
}

// Package.json script
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

////// Script //////
func main() {
    welcomeMessage()
    projectName := promptProjectName()
    initializeGit(projectName)
    createDirectories()

    if err := setupFrontend(); err != nil {
        fmt.Println("Error setting up frontend:", err)
        return
    }
    if err := setupBackend(); err != nil {
        fmt.Println("Error setting up backend:", err)
        return
    }
    if err := setupDatabase(); err != nil {
        fmt.Println("Error setting up database:", err)
        return
    }
    if err := setupOrm(); err != nil {
        fmt.Println("Error setting up ORM:", err)
        return
    }
    if err := updateEnvFile(); err != nil {
        fmt.Println("Error updating .env file:", err)
        return
    }
    if err := setupGitRemote(); err != nil {
        fmt.Println("Error setting up Git remote:", err)
        return
    }

    fmt.Println("Jeeeez! Project setup is ready, let's rip some code!")
}