.PHONY: clean

clean:
	@echo "Cleaning up processes..."
	-taskkill //F //IM t-log.exe 2>NUL || echo "No t-log.exe running."
	-taskkill //F //IM wails.exe 2>NUL || echo "No wails.exe running."
	@echo "Removing build artifacts..."
	-rm -rf build/bin/
	@echo "Done."

