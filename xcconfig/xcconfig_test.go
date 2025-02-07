package xcconfig

import (
	"errors"
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/bitrise-io/go-xcode/v2/xcconfig/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_WhenWritingXCConfigContent_ThenItShouldReturnFilePath(t *testing.T) {
	// Given
	var (
		testContent      = "TEST"
		testTempDir      = "temp_dir"
		expectedPath     = filepath.Join(testTempDir, "temp.xcconfig")
		mockPathModifier = new(mocks.PathModifier)
		mockPathChecker  = new(mocks.PathChecker)
		mockPathProvider = new(mocks.PathProvider)
		mockFileManager  = new(mocks.FileManager)
	)

	mockPathProvider.On("CreateTempDir", "").Return(testTempDir, nil)
	mockFileManager.On("Write", expectedPath, testContent, fs.FileMode(0644)).Return(nil)
	xcconfigWriter := NewWriter(mockPathProvider, mockFileManager, mockPathChecker, mockPathModifier)

	// When
	path, err := xcconfigWriter.Write(testContent)

	// Then
	if assert.NoError(t, err) {
		assert.Equal(t, expectedPath, path)
	}

	mockPathProvider.AssertExpectations(t)
	mockFileManager.AssertExpectations(t)
	mockPathModifier.AssertExpectations(t)
	mockPathChecker.AssertExpectations(t)
}

func Test_XCConfigInput_NonExistentPathErrors(t *testing.T) {
	// Given
	var (
		testContent      = "TEST.xcconfig"
		mockPathModifier = new(mocks.PathModifier)
		mockPathChecker  = new(mocks.PathChecker)
		mockPathProvider = new(mocks.PathProvider)
		mockFileManager  = new(mocks.FileManager)
	)

	mockPathModifier.On("AbsPath", testContent).Return(testContent, nil)
	mockPathChecker.On("IsPathExists", testContent).Return(false, errors.New("path does not exist"))
	xcconfigWriter := NewWriter(mockPathProvider, mockFileManager, mockPathChecker, mockPathModifier)

	// When
	path, err := xcconfigWriter.Write(testContent)

	// Then
	assert.Error(t, err)
	assert.Equal(t, path, "")

	mockPathProvider.AssertExpectations(t)
	mockFileManager.AssertExpectations(t)
	mockPathModifier.AssertExpectations(t)
	mockPathChecker.AssertExpectations(t)
}

func Test_XCConfigInput_CorrectInputPathReturnSamePath(t *testing.T) {
	// Given
	var (
		input            = "TEST.xcconfig"
		mockPathChecker  = new(mocks.PathChecker)
		mockPathProvider = new(mocks.PathProvider)
		mockFileManager  = new(mocks.FileManager)
		mockPathModifier = new(mocks.PathModifier)
	)

	mockPathModifier.On("AbsPath", input).Return(input, nil)
	mockPathChecker.On("IsPathExists", input).Return(true, nil)
	xcconfigWriter := NewWriter(mockPathProvider, mockFileManager, mockPathChecker, mockPathModifier)

	// When
	path, err := xcconfigWriter.Write(input)

	// Then
	if assert.NoError(t, err) {
		assert.Equal(t, path, input)
	}

	mockPathProvider.AssertExpectations(t)
	mockFileManager.AssertExpectations(t)
	mockPathModifier.AssertExpectations(t)
	mockPathChecker.AssertExpectations(t)
}
