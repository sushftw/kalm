import { createStyles, Theme } from "@material-ui/core";
import { WithStyles, withStyles } from "@material-ui/styles";
import React from "react";
import { connect } from "react-redux";
import { RootState } from "reducers";
import { TDispatch } from "types";
import { LEFT_SECTION_WIDTH } from "../pages/BasePage";
import { APP_BAR_HEIGHT } from "./AppBar";
import { H4 } from "../widgets/Label";

export const SECOND_HEADER_HEIGHT = 48;

const mapStateToProps = (state: RootState) => {
  return {};
};

const styles = (theme: Theme) =>
  createStyles({
    root: {
      position: "fixed",
      zIndex: 1201,
      top: APP_BAR_HEIGHT,
      height: SECOND_HEADER_HEIGHT,
      width: "100%",
      background: "#FFFFFF",
      borderBottom: "1px solid #EEEEEE",
      display: "flex"
    },
    left: {
      width: LEFT_SECTION_WIDTH,
      height: SECOND_HEADER_HEIGHT,
      borderRight: "1px solid #EEEEEE"
    },
    leftTextContainer: {
      display: "flex",
      alignItems: "center",
      paddingLeft: 32,
      borderRight: "1px solid #EEEEEE"
    },
    right: {
      flex: 1,
      height: SECOND_HEADER_HEIGHT
    }
  });

interface Props extends WithStyles<typeof styles>, ReturnType<typeof mapStateToProps> {
  dispatch: TDispatch;
  left: any;
  right?: any;
}

interface State {}

class SecondHeaderRaw extends React.PureComponent<Props, State> {
  constructor(props: Props) {
    super(props);

    this.state = {};
  }

  render() {
    const { classes, left, right } = this.props;

    return (
      <div className={classes.root}>
        <div className={`${classes.left} ${typeof left === "string" ? classes.leftTextContainer : ""}`}>
          {typeof left === "string" ? <H4>{left}</H4> : left}
        </div>
        <div className={classes.right}>{right}</div>
      </div>
    );
  }
}

export const SecondHeader = connect(mapStateToProps)(withStyles(styles)(SecondHeaderRaw));