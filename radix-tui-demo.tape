Output radix-tui-demo.gif

# Terminal appearance and environment
Set FontSize 24
Set FontFamily "MesloLGS NF"
Set Theme "Catppuccin Mocha"
Set Width 1800
Set Height 900
Set Padding 10
Set Margin 30
Set MarginFill "#674EFF"
Set BorderRadius 12
Set Shell zsh

# Build the app and set demo mode
Hide
Type "go build -o radix-tui ."
Enter
Type `export RADIX_TUI_DEMO=1`
Enter
Type "clear"
Enter
Show

# Launch the TUI with mocked data
Type "./radix-tui"
Enter
Sleep 2.0s

# Navigate applications list and open one
Down
Sleep 600ms
Up
Sleep 600ms
Enter
Sleep 1.2s

# Browse pipeline jobs (right pane)
Down
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 500ms
Up
Sleep 700ms

# Return to applications and pick another
Escape
Sleep 800ms
Down
Sleep 500ms
Down
Sleep 500ms
Enter
Sleep 1.2s

# Scroll pipeline again
Down
Sleep 500ms
Down
Sleep 500ms
Up
Sleep 600ms

# Quit the demo
Type "q"
Sleep 800ms