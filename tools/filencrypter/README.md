# FileCrypt

FileCrypt is a command-line tool to encrypt and decrypt image and video files in a specified folder using strong encryption.

## Usage

```sh
filecrypt <encrypt|decrypt> <folder> <password>

- <encrypt|decrypt>: Specify whether to encrypt or decrypt the files.
- <folder>: The folder containing the files to be encrypted or decrypted.
- <password>: The password used for encryption or decryption.
```

## Example

To encrypt all image and video files in the media folder with the password mysecretpassword:

``` sh
filecrypt encrypt media mysecretpassword
```

To decrypt all encrypted files in the media folder with the password mysecretpassword:

``` sh
filecrypt decrypt media mysecretpassword
```

## Supported File Types

- jpg
- png
- mp4
- mov
