import os


def main():
    for filename in os.listdir("pla"):
        if filename.endswith(".pla"):
            # Perform operations on the PLA file, e.g., convert to BLIF, etc.
            pla_file_path = os.path.join("pla", filename)
            blif_file_path = os.path.join("blif", filename[:-3] + "blif")

            # Invoke ABC to convert PLA to BLIF
            abc_command = f'abc -c "read_pla {pla_file_path}; strash; dc2; map -M 2; write_blif {blif_file_path}"'
            os.system(abc_command)


if __name__ == "__main__":
    main()
