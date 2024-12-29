# FileCrypt

FileCrypt is a command-line tool to encrypt and decrypt any file format in a specified folder using strong encryption. It maintains the subfolder structure and uses a subfolder-wise different counter for files.

## Features

- Encrypt and decrypt any file format and any size.
- Strong encryption using AES-256-CBC.
- Maintains the subfolder structure.

## Usage

```sh
filecrypt -m <mode/> -s <source_folder/> -d <destination_folder/> -p <password/>

-m: Mode of operation. Supported values are encrypt and decrypt. Default is encrypt.
-s: Source folder. The folder containing the files to be encrypted or decrypted.
-d: Destination folder. The folder where the encrypted or decrypted files will be saved. If not provided, it defaults to the source folder.
-p: Password. The password used for encryption or decryption. Uses the first 32 bytes of the password.
```

## Example

To encrypt all files in the `/path/to/original/folder` folder and save the encrypted files in the `/path/to/encrypted/folder` folder with the password `mysecretpassword`:

``` sh
filecrypt -m encrypt -s /path/to/original/folder -d /path/to/encrypted/folder -p mysecretpassword
```

To decrypt all files in the `/path/to/encrypted/folder` folder and save the decrypted files in the `/path/to/decrypted/folder` folder with the password `mysecretpassword`:

``` sh
filecrypt -m decrypt -s /path/to/encrypted/folder -d /path/to/decrypted/folder -p mysecretpassword
```

## File Naming Format

### Encrypted Files

encrypted files are named in the format

``` markdown
YYYYMMDD_format_xxxx.format.enc
```

- YYYYMMDD: The date when the file was created/updated.
- format: The file format (e.g., jpg, png, mp4).
- xxxx: A subfolder-wise counter for the files of that format.

### Decrypted Files

decrypted files are named in the format

``` markdown
YYYYMMDD_format_xxxx.format
```

- YYYYMMDD_format_xxxx.format: The original file name before decryption.
