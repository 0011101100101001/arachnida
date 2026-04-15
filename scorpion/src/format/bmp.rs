use std::io;
use std::path::Path;

#[derive(Default)]
pub struct Bmp {
    pub width: u32,
    pub height: i32, // Can be negative for top-down DIBs
    pub bits_per_pixel: u16,
    pub compression: u32,
    pub image_size: u32,
    pub colors_used: u32,
}

impl Bmp {
    pub fn read_metadata(&self, path: &Path) -> io::Result<String> {
        Ok(format!("BMP metadata for {}", path.display()))
    }
}
