package estructuras

type MBR struct {
	Mbr_tamano         int64
	Mbr_fecha_creacion [20]byte
	Mbr_dsk_signature  int64
	Dsk_fit            byte
	Mbr_partitions     [4]Partition
}

type Partition struct {
	Part_status byte
	Part_type   byte
	Part_fit    byte
	Part_start  int64
	Part_size   int64
	Part_name   [16]byte
}

type EBR struct {
	Part_status byte
	Part_fit    byte
	Part_start  int64
	Part_size   int64
	Part_next   int64
	Part_name   [16]byte
}

type SuperBloque struct {
	S_filesystem_type   int64
	S_inodes_count      int64
	S_blocks_count      int64
	S_free_blocks_count int64
	S_free_inodes_count int64
	S_mtime             [20]byte
	S_mnt_count         int64
	S_magic             int64
	S_inode_size        int64
	S_block_size        int64
	S_first_ino         int64
	S_first_blo         int64
	S_bm_inode_start    int64
	S_bm_block_start    int64
	S_inode_start       int64
	S_block_start       int64
}

type TablaInodo struct {
	I_uid   int64
	I_gid   int64
	I_size  int64
	I_atime [20]byte
	I_ctime [20]byte
	I_mtime [20]byte
	I_block [16]int64
	I_type  byte
	I_perm  int64
}

type BloqueDeCarpetas struct {
	B_content [4]Content
}

type Content struct {
	B_name  [12]byte
	B_inodo int32
}

type BloqueDeArchivos struct {
	B_content [64]byte
}
