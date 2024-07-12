mod bencode;
mod bitmap;
mod cli;
mod command;
mod peer;
mod torrent;
mod tracker;

use anyhow::Result;
use bencode::BencodeValue;
use clap::Parser;
use cli::Args;
use cli::Commands;
use torus::my_serde_bencode::from_str;

#[tokio::main]
async fn main() -> Result<()> {
    let args = Args::parse();

    match args.cmd {
        Commands::Decode { value } => {
            let bencode_value = from_str::<BencodeValue>(&value)?;
            println!("{bencode_value}");
        }
        Commands::Info { torrent_file } => command::info(&torrent_file)?,
        Commands::Peers { torrent_file } => command::peers(&torrent_file).await?,
        Commands::Handshake { torrent_file, ip } => command::handshake(&torrent_file).await?,
        Commands::DownloadPiece {
            output_file,
            torrent_file,
            piece_index,
        } => command::download_and_write_piece(&output_file, &torrent_file, piece_index).await?,
        Commands::Download {
            output_file,
            torrent_file,
        } => command::download(&output_file, &torrent_file).await?,
    }

    Ok(())
}
