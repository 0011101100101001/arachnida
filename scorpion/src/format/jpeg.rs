use std::io;
use std::path::Path;
use std::collections::HashMap;

#[derive(Default)]
pub struct Jpeg {
    pub width: u16,
    pub height: u16,
    pub color_components: u8, // e.g., 3 for RGB/YCbCr
    pub exif_data: HashMap<String, String>, // Camera Make, Date, GPS, etc.
}

impl Jpeg {
    pub fn read_metadata(&self, path: &Path) -> io::Result<String> {
        Ok(format!("JPEG metadata for {}", path.display()))
    }
}
