use std::io;
use std::path::Path;

#[derive(Default)]
pub struct Gif {
    pub version: String, // Usually "87a" or "89a"
    pub width: u16,
    pub height: u16,
    pub has_global_color_table: bool,
    pub background_color_index: u8,
    pub frame_count: u32,
    pub is_animated: bool,
}

impl Gif {
    pub fn read_metadata(&self, path: &Path) -> io::Result<String> {
        Ok(format!("GIF metadata for {}", path.display()))
    }
}
