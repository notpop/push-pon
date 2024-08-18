.PHONY: flash monitor

flash:
	tinygo flash --target waveshare-rp2040-zero --size short .

flashMonitor:
	tinygo flash --target waveshare-rp2040-zero --size short . && tinygo monitor

monitor:
	tinygo monitor

ports:
	tinygo ports
