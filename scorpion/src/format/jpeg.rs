use std::io;
use std::path::Path;

use super::Format;

pub struct Jpeg;

impl Format for Jpeg {
    fn extension(&self) -> &'static str {
        "jpeg"
    }

    fn read_metadata(&self, path: &Path) -> io::Result<String> {
        Ok(format!("JPEG metadata for {}", path.display()))
    }
}
