import { Button, createStyles, Theme, withStyles, WithStyles } from "@material-ui/core";
import React from "react";
import { connect } from "react-redux";
import { RootState } from "reducers";
import { TDispatchProp } from "types";
import { HttpRoute } from "types/route";
import { IconButtonWithTooltip } from "./IconButtonWithTooltip";
import { OpenInBrowserIcon } from "./Icon";
import { ClusterInfo } from "types/cluster";

const styles = (theme: Theme) =>
  createStyles({
    root: {},
  });

const mapStateToProps = (state: RootState) => {
  return {
    clusterInfo: state.get("cluster").get("info"),
  };
};

interface Props extends WithStyles<typeof styles>, ReturnType<typeof mapStateToProps>, TDispatchProp {
  route: HttpRoute;
  showIconButton?: boolean;
}

class OpenInBrowserRaw extends React.PureComponent<Props> {
  public render() {
    const { route, clusterInfo, showIconButton } = this.props;
    if (showIconButton) {
      return (
        <IconButtonWithTooltip
          tooltipTitle="Open In Browser"
          disabled={!route.get("methods").includes("GET")}
          href={getRouteUrl(route, clusterInfo)}
          // @ts-ignore
          target="_blank"
          rel="noreferrer"
        >
          <OpenInBrowserIcon />
        </IconButtonWithTooltip>
      );
    }

    return (
      <Button
        size="small"
        variant="outlined"
        disabled={!route.get("methods").includes("GET")}
        href={getRouteUrl(route, clusterInfo)}
        target="_blank"
        rel="noreferrer"
      >
        Open in browser
      </Button>
    );
  }
}

export const OpenInBrowser = withStyles(styles)(connect(mapStateToProps)(OpenInBrowserRaw));

/**
 * Get the URL string of a Route object
 */
export const getRouteUrl = (route: HttpRoute, clusterInfo: ClusterInfo, customHost?: string) => {
  let host = customHost ? customHost : route.get("hosts").first("*");
  const scheme = route.get("schemes").first();
  // TODO: Grabbing the first host and first url is probably not sufficient here.
  // What is the right behavior when we are dealing with a row with 2 hosts and 2 paths?
  const path = route.get("paths").first("/");

  if (host === "*") {
    host =
      (clusterInfo.get("ingressIP") || clusterInfo.get("ingressHostname")) +
      ":" +
      clusterInfo.get(scheme === "https" ? "httpsPort" : "httpPort");
  }

  if (host.includes("*")) {
    host = host.replace("*", "wildcard");
  }

  return scheme + "://" + host + path;
};
