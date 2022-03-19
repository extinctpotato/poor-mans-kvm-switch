ARD 		= arduino-cli
FQBN 		= arduino:avr:leonardo

BUILDDIR 	= build
SRCDIR 		= .

SKETCH		= ard_serial_hid_controller

compile: $(BUILDDIR)/$(SKETCH).ino.with_bootloader.hex

$(BUILDDIR)/%.ino.with_bootloader.hex: %.ino
	@echo "Compiling $< into $@"
	$(ARD) compile -b $(FQBN) --build-path $(BUILDDIR) $<
