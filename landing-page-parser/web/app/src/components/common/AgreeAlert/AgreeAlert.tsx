import * as React from 'react';
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import {useEffect} from "react";

const AlertDialog: React.FC<{
    title: string,
    text: string,
    onAgree:()=>void
    onClose:()=>void
    open: boolean
}> = (props) => {
    const [open, setOpen] = React.useState(props.open);

    useEffect(() => {
        setOpen(props.open)
    }, [props.open]);

    return (
        <React.Fragment>
            <Dialog
                open={open}
                onClose={props.onClose}
                aria-labelledby="alert-dialog-title"
                aria-describedby="alert-dialog-description"
            >
                <DialogTitle id="alert-dialog-title">
                    {props.title}
                </DialogTitle>
                <DialogContent>
                    <DialogContentText id="alert-dialog-description">
                        {props.text}
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={props.onClose}>Отмена</Button>
                    <Button onClick={()=>{
                        props.onClose()
                        props.onAgree()
                    }} autoFocus>
                        Подтвердить
                    </Button>
                </DialogActions>
            </Dialog>
        </React.Fragment>
    );
}

export default AlertDialog
