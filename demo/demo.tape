# Terminal demo for the README and the website.
# Regenerate with: go build -o demo/bin/apg . && vhs demo/demo.tape
# demo/pbcopy is a no-op stand-in so --copy leaves your clipboard alone.

Output docs/demo.gif

Set Shell bash
Set FontSize 18
Set Width 1000
Set Height 470
Set Padding 24
Set Theme "Catppuccin Mocha"
Set TypingSpeed 45ms

Hide
Type "export PATH=$PWD/demo/bin:$PWD/demo:$PATH PS1='$ ' && clear"
Enter
Show

Type "apg"
Sleep 300ms
Enter
Sleep 1.5s

Type "apg -n 3 --alnum --length 20"
Sleep 300ms
Enter
Sleep 1.5s

Type "apg -a 0 -n 3 -m 8 -x 10 -M sncl -t"
Sleep 300ms
Enter
Sleep 2s

Type "apg --pin --length 6"
Sleep 300ms
Enter
Sleep 1.5s

Type "apg -E 'O0Il1|' --length 16 --copy"
Sleep 300ms
Enter
Sleep 3s
