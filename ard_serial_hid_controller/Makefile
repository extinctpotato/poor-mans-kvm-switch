ARD 		= arduino-cli
BOARD_FQBN 	= arduino:avr:leonardo
BOARD_TTY	= /dev/ttyACM0

BUILDDIR 	= build
LIBDIR		= lib

SKETCH		= ard_serial_hid_controller
SKETCHLIBS	= $(addprefix --libraries ,$(wildcard $(LIBDIR)/*))
SKETCHHEX	= $(BUILDDIR)/$(SKETCH).ino.with_bootloader.hex

build: $(SKETCHHEX)

upload: build
	$(ARD) upload -b $(BOARD_FQBN) -p $(BOARD_TTY) --input-dir=$(BUILDDIR)

$(BUILDDIR)/%.ino.with_bootloader.hex: %.ino
	@echo "Compiling $< into $@"
	$(ARD) compile -b $(BOARD_FQBN) $(SKETCHLIBS) --build-path $(BUILDDIR) $<
