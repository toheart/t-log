.PHONY: clean

clean:
	@echo "Cleaning up processes..."
	-taskkill //F //IM t-log-dev.exe 2>NUL || echo "No t-log-dev.exe running."
	-taskkill //F //IM wails.exe 2>NUL || echo "No wails.exe running."
	@echo "Removing build artifacts..."
	-rm -rf build/bin/
	@echo "Done."

