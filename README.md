```
===============================================================================
                          CHIP-8 EMULATOR IN GO
===============================================================================

Description:
------------
This is a chip-8 emulator built in go with implementations from scratch. Chip-8
was created in the mid-1970s by Joseph Weisbecker as an interpreted programming
language for early microcomputers like the COSMAC VIP and Telmac 1800. It was
designed to make game development easier by providing simple graphics and sound
capabilities through a virtual machine.

This implementation is fully functional and is powered by ebitengine which is a
game engine for game development in golang. I am using this engine to play aud-
io and display.

Features:
---------
- Fully implements Chip-8 instructions keeping the implementation as true to the
  original as possible
- Should be able to run any Chip-8 game compatible with other well known implem-
  entations
```
![Rock-Paper-Scissors-ROM](https://github.com/user-attachments/assets/2a1c8a6d-c787-4ad4-84ff-9e6cac469cab)

```
Requirements:
-------------
* Needs to be run through a terminal, since other UI elements are not added yet

Installation:
-------------
1. Step 1:
   Install the stable GitHub Release available.
        OR
   Install through go cli (requires go to be installed)
   --------------------------------------------------------------
   RUN:   go install github.com/Shreyaskr1409/go-chip8-emu@latest
   --------------------------------------------------------------
   This will install the executable in you Go "bin" directory (typically ~/go/bin
   on Unix/macOS or %USERPROFILE%\go\bin on Windows).

2. Step 2:
   Locate the build file and execute it as:
   ------------------------------------
   RUN:   go-chip8-emu /path/to/rom.ch8
   ------------------------------------
   REPLACE /path/to/rom.ch8 with you ROM file.

Configuration:
---------------
Configutation is yet to be added, will be added soon enough. It will enable sca-
ling of the game window.

License:
--------
This project is licensed under the GPL-3.0 license - see the LICENSE file for
details.

PROJECT:
--------
Project URL: https://github.com/Shreyaskr1409/go-chip8-emu

===============================================================================
                            END OF README
===============================================================================
```
