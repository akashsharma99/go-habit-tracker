# Habit Tracker (Terminal)

A simple terminal-based habit tracker built in Go using the Bubble Tea TUI framework. This app helps you build and maintain good habits by tracking your daily progress in a clean, interactive terminal UI.

## Features
- Add new habits
- Mark habits as completed for today
- View completion rate for each habit
- Navigate habits with arrow keys
- Data is stored locally

## Installation

1. **Clone the repository:**
   ```sh
   git clone https://github.com/akashsharma99/go-habit-tracker.git
   cd go-habit-tracker
   ```
2. **Install dependencies:**
   ```sh
   go mod tidy
   ```
3. **Build the app:**
   ```sh
   go build -o habit-tracker
   ```

## Usage

Run the app from your terminal:

```sh
./habit-tracker
```

### Controls
- **Up/Down Arrow:** Move cursor between habits
- **Space/Enter:** Toggle completion status for today
- **a:** Add a new habit
- **q:** Quit the app
- **esc:** Cancel adding a new habit

### Data Storage
- Habit data is stored locally using sqlite db and is not synced to the cloud.
- Log files are created in your home directory under `.habit-tracker/app.log` for debugging.

## Contributing
Pull requests and suggestions are welcome!
