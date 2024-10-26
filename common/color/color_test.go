package color

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Red(t *testing.T) {
	reader, writer, err := os.Pipe()
	require.Nil(t, err)

	originalStdout := os.Stdout
	defer func() {
		os.Stdout = originalStdout // Restore original stdout after test
		reader.Close()
		writer.Close()
	}()

	os.Stdout = writer

	fmt.Print(Red("Hello, Red World!"))

	writer.Close() // Close writer to signal we're done writing
	var buf bytes.Buffer
	_, err = io.Copy(&buf, reader)
	require.Nil(t, err)

	expected := "\033[31mHello, Red World!\033[0m"
	require.Equal(t, expected, buf.String())
}

func Test_Green(t *testing.T) {
	reader, writer, err := os.Pipe()
	require.Nil(t, err)

	originalStdout := os.Stdout
	defer func() {
		os.Stdout = originalStdout // Restore original stdout after test
		reader.Close()
		writer.Close()
	}()

	os.Stdout = writer

	fmt.Print(Green("Hello, Red World!"))

	writer.Close() // Close writer to signal we're done writing
	var buf bytes.Buffer
	_, err = io.Copy(&buf, reader)
	require.Nil(t, err)

	expected := "\033[32mHello, Red World!\033[0m"
	require.Equal(t, expected, buf.String())
}

func Test_Yellow(t *testing.T) {
	reader, writer, err := os.Pipe()
	require.Nil(t, err)

	originalStdout := os.Stdout
	defer func() {
		os.Stdout = originalStdout // Restore original stdout after test
		reader.Close()
		writer.Close()
	}()

	os.Stdout = writer

	fmt.Print(Yellow("Hello, Red World!"))

	writer.Close() // Close writer to signal we're done writing
	var buf bytes.Buffer
	_, err = io.Copy(&buf, reader)
	require.Nil(t, err)

	expected := "\033[33mHello, Red World!\033[0m"
	require.Equal(t, expected, buf.String())
}

func Test_Blue(t *testing.T) {
	reader, writer, err := os.Pipe()
	require.Nil(t, err)

	originalStdout := os.Stdout
	defer func() {
		os.Stdout = originalStdout // Restore original stdout after test
		reader.Close()
		writer.Close()
	}()

	os.Stdout = writer

	fmt.Print(Blue("Hello, Red World!"))

	writer.Close() // Close writer to signal we're done writing
	var buf bytes.Buffer
	_, err = io.Copy(&buf, reader)
	require.Nil(t, err)

	expected := "\033[34mHello, Red World!\033[0m"
	require.Equal(t, expected, buf.String())
}

func Test_Magenta(t *testing.T) {
	reader, writer, err := os.Pipe()
	require.Nil(t, err)

	originalStdout := os.Stdout
	defer func() {
		os.Stdout = originalStdout // Restore original stdout after test
		reader.Close()
		writer.Close()
	}()

	os.Stdout = writer

	fmt.Print(Magenta("Hello, Red World!"))

	writer.Close() // Close writer to signal we're done writing
	var buf bytes.Buffer
	_, err = io.Copy(&buf, reader)
	require.Nil(t, err)

	expected := "\033[35mHello, Red World!\033[0m"
	require.Equal(t, expected, buf.String())
}

func Test_Cyan(t *testing.T) {
	reader, writer, err := os.Pipe()
	require.Nil(t, err)

	originalStdout := os.Stdout
	defer func() {
		os.Stdout = originalStdout // Restore original stdout after test
		reader.Close()
		writer.Close()
	}()

	os.Stdout = writer

	fmt.Print(Cyan("Hello, Red World!"))

	writer.Close() // Close writer to signal we're done writing
	var buf bytes.Buffer
	_, err = io.Copy(&buf, reader)
	require.Nil(t, err)

	expected := "\033[36mHello, Red World!\033[0m"
	require.Equal(t, expected, buf.String())
}

func Test_Gray(t *testing.T) {
	reader, writer, err := os.Pipe()
	require.Nil(t, err)

	originalStdout := os.Stdout
	defer func() {
		os.Stdout = originalStdout // Restore original stdout after test
		reader.Close()
		writer.Close()
	}()

	os.Stdout = writer

	fmt.Print(Gray("Hello, Red World!"))

	writer.Close() // Close writer to signal we're done writing
	var buf bytes.Buffer
	_, err = io.Copy(&buf, reader)
	require.Nil(t, err)

	expected := "\033[37mHello, Red World!\033[0m"
	require.Equal(t, expected, buf.String())
}

func Test_White(t *testing.T) {
	reader, writer, err := os.Pipe()
	require.Nil(t, err)

	originalStdout := os.Stdout
	defer func() {
		os.Stdout = originalStdout // Restore original stdout after test
		reader.Close()
		writer.Close()
	}()

	os.Stdout = writer

	fmt.Print(White("Hello, Red World!"))

	writer.Close() // Close writer to signal we're done writing
	var buf bytes.Buffer
	_, err = io.Copy(&buf, reader)
	require.Nil(t, err)

	expected := "\033[97mHello, Red World!\033[0m"
	require.Equal(t, expected, buf.String())
}
