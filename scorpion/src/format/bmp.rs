use std::io;
use std::path::Path;

use super::Format;

pub struct Bmp;

impl Format for Bmp {
    fn extension(&self) -> &'static str {
        "bmp"
    }

    fn read_metadata(&self, path: &Path) -> io::Result<String> {
        Ok(format!("BMP metadata for {}", path.display()))
    }
}
