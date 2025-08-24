// TODO:  Define your log entry structure (key + value as bytes)
// TODO: Implement escaping/unescaping functions for § and newlines
// TODO: Write functions to serialize entries to strings and parse them back
// TODO: Basic unit tests for edge cases (empty keys, values with delimiters, etc.)

// TODO: File writing: append entries to log file
// TODO: File reading: read and parse entries from log file
// TODO: Handle file creation, opening existing files
// TODO: Error handling for I/O operations

use std::fs::{File, OpenOptions};
use std::io::{ErrorKind, Read, Write};

const DB_FILE: &str = "db/db.txt";

struct Entry {
    key: Vec<u8>,
    value: Vec<u8>,
}

fn put(entry: &Entry, db: &mut File) -> std::io::Result<()> {
    let mut data: Vec<u8> = Vec::new();
    data.extend_from_slice(&(entry.key.len() as u32).to_le_bytes());
    data.extend_from_slice(&(entry.value.len() as u32).to_le_bytes());
    data.extend_from_slice(&entry.key);
    data.extend_from_slice(&entry.value);
    db.write_all(&data)?;
    Ok(())
}

fn read_entry(db: &mut File) -> std::io::Result<Option<Entry>> {
    let mut len_buf = [0u8; 4];
    match db.read_exact(&mut len_buf) {
        Ok(()) => {},
        Err(e) if e.kind() == std::io::ErrorKind::UnexpectedEof => return Ok(None),
        Err(e) => return Err(e),
    }
    let key_len = u32::from_le_bytes(len_buf) as usize;
    db.read_exact(&mut len_buf)?;
    let value_len = u32::from_le_bytes(len_buf) as usize;
    let mut key = vec![0u8; key_len];
    db.read_exact(&mut key)?;
    let mut value = vec![0u8; value_len];
    db.read_exact(&mut value)?;

    Ok(Some(Entry { key, value }))
}

fn get(key: &Vec<u8>, db: &mut File) -> std::io::Result<Option<Entry>> {
    while let Some(entry) = read_entry(db)? {
        if entry.key == *key {
            return Ok(Some(entry));
        }
    }
    Ok(None)
}


fn main() -> std::io::Result<()> {
    let mut open_options = OpenOptions::new();
    let mut db_write_handle = open_options.append(true).read(true).open(DB_FILE).unwrap_or_else(|error| match error.kind() {
        ErrorKind::NotFound => {
            match open_options.append(true).read(true).create(true).open(DB_FILE) {
                Ok(created) => { created }
                Err(error) => {
                    panic!("Problem creating file: {error:?}");
                }
            }
        },
        _ => {
            panic!("Problem opening file: {error:?}");
        }
    });
    let mut db_read_handle = open_options.read(true).open(DB_FILE).unwrap_or_else(|error| match error.kind() {
        ErrorKind::NotFound => {
            match open_options.read(true).create(true).open(DB_FILE) {
                Ok(created) => { created }
                Err(error) => {
                    panic!("Problem creating file: {error:?}");
                }
            }
        },
        _ => {
            panic!("Problem opening file: {error:?}");
        }
    });

    let entry = Entry {key: Vec::from("test".as_bytes()), value: Vec::from("hello world".as_bytes())};

    // put an entry
    // - create entry
    // - serialize key§value
    // - append to file
    put(&entry, &mut db_write_handle)?;
    // get a value
    // - search entry by key
    // - deserialize entry
    // - return entry deserialized
    match get(&entry.key, &mut db_read_handle) {
        Ok(Some(entry)) => {
            println!("found")
        }
        Ok(None) => {
            println!("entry not found");
        }
        Err(error) => {
            println!("error: {error:?}");
        }
    }

    Ok(())

}
