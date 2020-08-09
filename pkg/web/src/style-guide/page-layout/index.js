// @flow

import * as React from "react";
import { Unstable_AppNavBar as AppNavBar } from "baseui/app-nav-bar";
import { useStyletron } from "baseui";
import { Cell, Grid } from "baseui/layout-grid";
import { Link, useLocation } from "react-router-dom";
import LogoSvg from "../illustrations/logo";

type Props = {
  children?: React.Node,
};

type NavLinkProps = {
  label: string,
  pathname: string,
};

const NavLink = (props: NavLinkProps) => {
  const [css] = useStyletron();

  return (
    <Link
      to={props.pathname}
      className={css({
        textDecoration: "none",
        color: "inherit",
        ":hover": { color: "inherit" },
        ":visited": { color: "inherit" },
      })}
    >
      {props.label}
    </Link>
  );
};

const PageLayout = (props: Props) => {
  const location = useLocation();
  const [css] = useStyletron();

  return (
    <>
      <div className={css({ position: "fixed", top: 0, left: 0, right: 0 })}>
        <AppNavBar
          appDisplayName={
            <Link to={"/"}>
              <LogoSvg height={32} />
            </Link>
          }
          mainNav={[
            {
              item: { label: "Audience", pathname: "/audience" },
              mapItemToString: (item) => item.label,
              mapItemToNode: (item) => <NavLink {...item} />,
            },
            {
              item: { label: "Templates", pathname: "/templates" },
              mapItemToString: (item) => item.label,
              mapItemToNode: (item) => <NavLink {...item} />,
            },
            {
              item: { label: "Campaigns", pathname: "/campaigns" },
              mapItemToString: (item) => item.label,
              mapItemToNode: (item) => <NavLink {...item} />,
            },
          ]}
          isNavItemActive={({ item: navItem }) =>
            location.pathname.startsWith(navItem.item.pathname)
          }
          onNavItemSelect={() => {}}
          username="DropList User"
          userNav={[
            {
              item: { label: "Credits", pathname: "/credits" },
              mapItemToString: (item) => item.label,
              mapItemToNode: (item) => <NavLink {...item} />,
            },
            {
              item: { label: "AWS Settings", pathname: "/profile" },
              mapItemToString: (item) => item.label,
              mapItemToNode: (item) => <NavLink {...item} />,
            },
            {
              item: { label: "Logout", pathname: "/logout" },
              mapItemToString: (item) => item.label,
              mapItemToNode: (item) => <NavLink {...item} />,
            },
          ]}
        />
      </div>

      <Grid
        overrides={{
          Grid: { style: { height: "100%", paddingTop: "112px" } },
        }}
      >
        <Cell span={12}>{props.children}</Cell>
      </Grid>
    </>
  );
};

export default PageLayout;
