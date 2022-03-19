ARD 		= arduino-cli
BOARD_FQBN 	= arduino:avr:leonardo
BOARD_TTY	= /dev/ttyACM0

BUILDDIR 	= build
SRCDIR 		= .

SKETCH		= ard_serial_hid_controller
SKETCHHEX	= $(BUILDDIR)/$(SKETCH).ino.with_bootloader.hex

build: $(SKETCHHEX)

upload: build
	$(ARD) upload -b $(BOARD_FQBN) -p $(BOARD_TTY) --input-dir=$(BUILDDIR)

$(BUILDDIR)/%.ino.with_bootloader.hex: %.ino
	@echo "Compiling $< into $@"
	$(ARD) compile -b $(BOARD_FQBN) --build-path $(BUILDDIR) $<
