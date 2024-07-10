use serde::{Deserialize, Serialize};
use sha1::{Digest, Sha1};
use torus::my_serde_bencode::to_bytes;

#[derive(Debug, Serialize, Deserialize)]
pub struct Torrent {
    pub announce: String,
    pub info: Info,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct Info {
    pub length: u32,
    pub name: String,
    #[serde(rename = "piece length")]
    pub piece_length: u32,
    #[serde(with = "serde_bytes")]
    pub pieces: Vec<u8>,
}

impl Torrent {
    pub fn info_hash(&self) -> Vec<u8> {
        let info = to_bytes(&self.info).unwrap();

        let mut hasher = Sha1::new();
        hasher.update(&info);
        hasher.finalize().to_vec()
    }
}
