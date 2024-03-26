build-elf:
	@cd lib/src/program && cargo prove build

cp-elf:
	@cp lib/src/program/elf/riscv32im-succinct-zkvm-elf lib/elf/riscv32im-succinct-zkvm-elf

build-prover:
	@cd lib/src/prover && cargo build --release

cp-prover:
	@cp lib/src/prover/target/release/libprover.a lib/libprover.a

all: build-elf cp-elf build-prover cp-prover 