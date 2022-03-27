import React, { Component, Fragment } from 'react';
import { Grid, Form, Input, Field, Checkbox, Button, Message } from '@b-design/ui';
import { If } from 'tsx-control-statements/components';
import Empty from '../../../../components/Empty';
import { createProjectRoles, updateProjectRoles } from '../../../../api/project';
import type { PermPolicies } from '../../../../interface/permPolicies';
import { ProjectRoleBase } from '../../../../interface/project';
import { checkName } from '../../../../utils/common';
import Translation from '../../../../components/Translation';
import i18n from '../../../../i18n';
import './index.less';

type Props = {
  projectName: string;
  projectRoles: ProjectRoleBase[];
  activeRoleName: string;
  activeRoleItem: ProjectRoleBase;
  isCreateProjectRoles: boolean;
  projectPermPolicies: PermPolicies[];
  onCreate: (activeRoleItem: string) => void;
};

type State = {
  loading: boolean;
};

type CheckBoxBase = { name: string }[];
const { Group: CheckboxGroup } = Checkbox;

class ProjectPermPoliciesDetail extends Component<Props, State> {
  field: Field;
  constructor(props: Props) {
    super(props);
    this.field = new Field(this);
    this.state = {
      loading: false,
    };
  }

  componentWillReceiveProps(nextProps: Props) {
    const { isCreateProjectRoles } = nextProps;
    if (isCreateProjectRoles) {
      this.field.setValues({
        name: '',
        alias: '',
        permPolicies: [],
      });
    } else {
      if (nextProps.activeRoleItem !== this.props.activeRoleItem) {
        this.field.setValues({
          name: nextProps.activeRoleItem.name || '',
          alias: nextProps.activeRoleItem.alias || '',
          permPolicies: this.initPermPoliciesStatus(nextProps.activeRoleItem),
        });
      }
    }
  }

  initPermPoliciesStatus = (activeItem: ProjectRoleBase) => {
    if (activeItem) {
      return (activeItem.permPolicies || []).map((item: { name: string }) => item.name);
    } else {
      return [];
    }
  };

  transCheckBoxData = (data: CheckBoxBase) => {
    return (data || []).map((item) => {
      return {
        value: item.name,
        label: item.name,
      };
    });
  };

  listPermPolicies = () => {
    const { activeRoleItem, isCreateProjectRoles, projectPermPolicies } = this.props;
    if (isCreateProjectRoles) {
      return this.transCheckBoxData(projectPermPolicies);
    } else {
      if (activeRoleItem && activeRoleItem.permPolicies) {
        return this.transCheckBoxData(activeRoleItem.permPolicies);
      }
    }
  };

  onSubmit = () => {
    this.field.validate((error: any, values: any) => {
      if (error) {
        return;
      }
      const { isCreateProjectRoles, projectName, activeRoleName } = this.props;
      const { name, alias, permPolicies } = values;
      const queryData = {
        projectName,
        roleName: activeRoleName,
      };
      const param = {
        name,
        alias,
        permPolicies,
      };
      this.setState({ loading: true });
      if (isCreateProjectRoles) {
        createProjectRoles(queryData, param)
          .then((res: { name: string }) => {
            this.setState({ loading: false });
            if (res) {
              Message.success(<Translation>Create role success</Translation>);
              this.props.onCreate(res.name);
            }
          })
          .finally(() => {
            this.setState({ loading: false });
          });
      } else {
        updateProjectRoles(queryData, param)
          .then((res: { name: string }) => {
            this.setState({ loading: false });
            if (res) {
              Message.success(<Translation>Update role success</Translation>);
              this.props.onCreate(res.name);
            }
          })
          .finally(() => {
            this.setState({ loading: false });
          });
      }
    });
  };

  render() {
    const init = this.field.init;
    const { Row, Col } = Grid;
    const FormItem = Form.Item;
    const formItemLayout = {
      labelCol: {
        fixedSpan: 6,
      },
      wrapperCol: {
        span: 20,
      },
    };
    const { projectRoles, isCreateProjectRoles } = this.props;
    return (
      <Fragment>
        <If condition={projectRoles && projectRoles.length === 0 && !isCreateProjectRoles}>
          <div className="project-role-empty-wrapper">
            <Empty />
          </div>
        </If>
        <If condition={(projectRoles && projectRoles.length !== 0) || isCreateProjectRoles}>
          <div className="auth-list-wrapper">
            <Form {...formItemLayout} field={this.field} className="auth-form-content">
              <Row>
                <Col span={12} style={{ padding: '16px 16px 0 30px' }}>
                  <FormItem
                    label={<Translation>Name</Translation>}
                    labelAlign="left"
                    required
                    className="font-weight-400"
                  >
                    <Input
                      name="name"
                      placeholder={i18n.t('Please enter').toString()}
                      maxLength={32}
                      disabled={isCreateProjectRoles ? false : true}
                      {...init('name', {
                        rules: [
                          {
                            required: true,
                            pattern: checkName,
                            message: <Translation>Please enter a valid name</Translation>,
                          },
                        ],
                      })}
                    />
                  </FormItem>
                </Col>
                <Col span={12} style={{ padding: '16px 16px 0 30px' }}>
                  <FormItem
                    label={<Translation>Alias</Translation>}
                    labelAlign="left"
                    className="font-weight-400"
                  >
                    <Input
                      name="alias"
                      placeholder={i18n.t('Please enter').toString()}
                      {...init('alias', {
                        rules: [
                          {
                            minLength: 2,
                            maxLength: 64,
                            message: 'Enter a string of 2 to 64 characters.',
                          },
                        ],
                      })}
                    />
                  </FormItem>
                </Col>
              </Row>
              <Row>
                <Col span={24} style={{ padding: '0 16px 16px 30px' }}>
                  <FormItem
                    label={<Translation>PermPolicies</Translation>}
                    labelAlign="left"
                    className="font-weight-400 permPolicies-wrapper"
                    required={true}
                  >
                    <CheckboxGroup
                      dataSource={this.listPermPolicies()}
                      {...init('permPolicies', {
                        rules: [
                          {
                            required: true,
                            type: 'array',
                            message: 'Choose one permPolicy',
                          },
                        ],
                      })}
                    />
                  </FormItem>
                </Col>
              </Row>
            </Form>
            <Button className="create-auth-btn" type="primary" onClick={this.onSubmit}>
              <Translation>{isCreateProjectRoles ? 'Create' : 'Update'}</Translation>
            </Button>
          </div>
        </If>
      </Fragment>
    );
  }
}

export default ProjectPermPoliciesDetail;