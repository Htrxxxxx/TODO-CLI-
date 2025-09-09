### TODO CLI

A simple command-line To-Do application built with Go.
This project helps you manage your tasks directly from the terminal (add, list, mark as done, delete, etc.).

### Features

Add tasks to your todo list

List all tasks

Mark tasks as completed

Delete tasks

Data persistence (tasks saved locally in a JSON file)

### Getting Started
1. Clone the repository
git clone https://github.com/<your-username>/todo-cli.git
cd todo-cli

2. Build the project
go build -o todo.exe


This will generate an executable file todo.exe (on Windows).
On Linux/Mac it will just be todo.

3. Run the app

You can now use the CLI commands:

./todo.exe add "Buy milk"
./todo.exe add "Finish homework"
./todo.exe list
./todo.exe done 1
./todo.exe delete 2


Example output:

1. [ ] Buy milk
2. [ ] Finish homework


After marking task 1 as done:

1. [x] Buy milk
2. [ ] Finish homework

### Commands
Command	Description
add "task name"	Add a new task
list	Show all tasks
done <task-id>	Mark task as completed
delete <task-id>	Delete a task
📂 Project Structure
TODO/
 ├── go.mod       # Go module file
 ├── TODO.go      # Main Go source code
 ├── todo.exe     # Compiled binary (ignored in .gitignore)

💡 Notes

On Windows Git Bash, you may need to run:

./todo.exe list


instead of todo.exe list.
