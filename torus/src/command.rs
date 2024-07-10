use std::fs;

use anyhow::{anyhow, Result};
use torus::my_serde_bencode::from_bytes;

use crate::{torrent::Torrent, tracker::discover_peers};

pub fn info(torrent_file: &str) -> Result<()> {
    let torrent = parse_torrent(torrent_file)?;
    let info_hash = torrent.info_hash();

    println!("Tracker URL: {}", torrent.announce);
    println!("Length: {}", torrent.info.length);
    println!("Info Hash: {}", hex::encode(info_hash));
    println!("Piece Length: {}", torrent.info.piece_length);
    println!("Piece Hashes:");
    for piece in torrent.info.pieces.chunks(20) {
        println!("{}", hex::encode(piece));
    }

    Ok(())
}

fn parse_torrent(torrent_file: &str) -> Result<Torrent> {
    let content = fs::read(torrent_file)?;
    let torrent = from_bytes::<Torrent>(&content)?;
    Ok(torrent)
}

pub async fn peers(torrent_file: &str) -> Result<()> {
    let torrent = parse_torrent(torrent_file)?;
    let peers = discover_peers(&torrent).await?;
    peers.iter().for_each(|peer| println!("{}", peer.to_url()));

    Ok(())
}

pub async fn handshake(torrent_file: &str) -> Result<()> {
    let torrent = parse_torrent(torrent_file)?;
    let peers = discover_peers(&torrent).await?;
    let peer = peers.first().ok_or(anyhow!("No peers found"))?;
    let peer = peer.connect(torrent.info_hash()).await?;
    println!("Peer ID: {}", peer.peer_id);

    Ok(())
}
