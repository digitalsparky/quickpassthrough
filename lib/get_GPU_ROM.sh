function get_GPU_ROM () {
    clear
    # Get the config paths
    source "$SCRIPTDIR/lib/paths.sh"

    VBIOS_PATH=$(find /sys/devices -name rom | grep $1)
    printf "We will now attempt to dump the vbios of your selected GPU.
Passing a VBIOS rom to the card used for passthrough is required for some cards, but not all.
Some cards also requires you to patch your VBIOS romfile, check online if this is neccessary for your card.
The VBIOS will be read from $VBIOS_PATH
This process will require the use of sudo and will run the following commands:

echo 1 | sudo tee $VBIOS_PATH
cat $VBIOS_PATH > $SCRIPTDIR/$QUICKEMU/vfio_card.rom
echo 0 | sudo tee $VBIOS_PATH

"
    read "Do you want to dump the VBIOS, choosing N will skip this step [y/N]: " YESNO
    case "${YESNO}" in
    [Yy]*)
        echo 1 | sudo tee $VBIOS_PATH
        cat $VBIOS_PATH > "$SCRIPTDIR/$QUICKEMU/vfio_card.rom"
        ROM_MD5=$(md5sum $VBIOS_PATH)
        echo 0 | sudo tee $VBIOS_PATH
        ROMFILE_MD5=$(md5sum "$SCRIPTDIR/$QUICKEMU/vfio_card.rom")

        if [ -f $SCRIPTDIR/$QUICKEMU/vfio_card.rom ];
        then
            if [ "$ROM_MD5" == "$ROMFILE_MD5" ];
            then
                echo "Checksums match!"
                echo "Dumping of VBIOS successful!"
                echo 'GPU_ROMFILE="vfio_card.rom"' >> "$SCRIPTDIR/$QUICKEMU/qemu-vfio_vars.conf"

                read "Press ENTER to continue." ENTER
            else
                echo "Checksums does not match!"
                echo "Dumping of VBIOS failed, skipping romfile"
                mv "$SCRIPTDIR/$QUICKEMU/vfio_card.rom" "$SCRIPTDIR/$QUICKEMU/vfio_card.rom.fail"
                echo 'GPU_ROMFILE=""' >> "$SCRIPTDIR/$QUICKEMU/qemu-vfio_vars.conf"

                read "Press ENTER to continue." ENTER
            fi
        else
            echo 'GPU_ROMFILE=""' >> "$SCRIPTDIR/$QUICKEMU/qemu-vfio_vars.conf"
        fi
    ;;
    [Nn]*)
        echo 'GPU_ROMFILE=""' >> "$SCRIPTDIR/$QUICKEMU/qemu-vfio_vars.conf"
    ;;
    *)
        echo 'GPU_ROMFILE=""' >> "$SCRIPTDIR/$QUICKEMU/qemu-vfio_vars.conf"
    ;;
    esac
}


function main () {
    SCRIPTDIR=$(dirname `which $0` | perl -pe "s/\/\.\.\/lib//")
    SCRIPTDIR="$SCRIPTDIR/.."
    
    get_GPU_ROM $1
}

main $1