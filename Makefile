APP_PKG ?= ./cmd/game
BUILD_DIR ?= build

.PHONY: mobile-android mobile-ios mobile-clean

mobile-android:
	@mkdir -p $(BUILD_DIR)
	ebitenmobile bind -target=android -o $(BUILD_DIR)/koro.aar $(APP_PKG)

mobile-ios:
	@mkdir -p $(BUILD_DIR)
	ebitenmobile bind -target=ios -o $(BUILD_DIR)/Koro.framework $(APP_PKG)

mobile-clean:
	rm -rf $(BUILD_DIR)
