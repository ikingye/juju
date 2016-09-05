#!/usr/bin/env python
"""Simple performance and scale tests."""

from __future__ import print_function

import argparse
from collections import defaultdict, namedtuple
from datetime import datetime
import logging
import os
import sys
import subprocess
import time

import rrdtool
from jinja2 import Template

from deploy_stack import (
    BootstrapManager,
)
from logbreakdown import (
    _render_ds_string,
    breakdown_log_by_timeframes,
)
import perf_graphing
from utility import (
    add_basic_testing_arguments,
    configure_logging,
)


__metaclass__ = type


class TimingData:

    strf_format = '%F %H:%M:%S'

    # Log breakdown uses the start/end too. Perhaps have a property for string
    # rep and a ds for the datetime.
    def __init__(self, start, end):
        self.start = start.strftime(self.strf_format)
        self.end = end.strftime(self.strf_format)
        self.seconds = int((end - start).total_seconds())

DeployDetails = namedtuple(
    'DeployDetails', ['name', 'applications', 'timings'])


SETUP_SCRIPT_PATH = 'perf_static/setup-perf-monitoring.sh'
COLLECTD_CONFIG_PATH = 'perf_static/collectd.conf'


log = logging.getLogger("assess_perf_test_simple")


def assess_perf_test_simple(bs_manager, upload_tools):
    # XXX
    # Find the actual cause for this!! (Something to do with the template and
    # the logs.)
    import sys
    reload(sys)  # NOQA
    sys.setdefaultencoding('utf-8')
    # XXX

    bs_start = datetime.utcnow()
    with bs_manager.booted_context(upload_tools):
        client = bs_manager.client
        admin_client = client.get_admin_client()
        admin_client.wait_for_started()
        bs_end = datetime.utcnow()
        try:
            apply_any_workarounds(client)
            bootstrap_timing = TimingData(bs_start, bs_end)

            setup_system_monitoring(admin_client)

            deploy_details = assess_deployment_perf(client)
        finally:
            results_dir = os.path.join(
                os.path.abspath(bs_manager.log_dir), 'performance_results/')
            os.makedirs(results_dir)
            admin_client.juju(
                'scp',
                ('--', '-r', '0:/var/lib/collectd/rrd/localhost/*',
                 results_dir)
            )

            try:
                admin_client.juju(
                    'scp', ('0:/tmp/mongodb-stats.log', results_dir)
                )
            except subprocess.CalledProcessError as e:
                log.error('Failed to copy mongodb stats: {}'.format(e))
            cleanup_start = datetime.utcnow()
    # Cleanup happens when we move out of context
    cleanup_end = datetime.utcnow()
    cleanup_timing = TimingData(cleanup_start, cleanup_end)
    deployments = dict(
        bootstrap=bootstrap_timing,
        deploys=[deploy_details],
        cleanup=cleanup_timing,
    )

    # Could be smarter about this.
    controller_log_file = os.path.join(
        bs_manager.log_dir,
        'controller',
        'machine-0',
        'machine-0.log.gz')

    generate_reports(controller_log_file, results_dir, deployments)


def apply_any_workarounds(client):
    # Work around mysql charm wanting 80% memory.
    if client.env.get_cloud() == 'lxd':
        constraint_cmd = [
            'lxc',
            'profile',
            'set',
            'juju-{}'.format(client.env.environment),
            'limits.memory',
            '2GB'
        ]
        subprocess.check_output(constraint_cmd)


def generate_reports(controller_log, results_dir, deployments):
    """Generate reports and graphs from run results."""
    cpu_image = generate_cpu_graph_image(results_dir)
    memory_image = generate_memory_graph_image(results_dir)
    network_image = generate_network_graph_image(results_dir)
    mongo_image = generate_mongo_graph_image(results_dir)

    log_message_chunks = get_log_message_in_timed_chunks(
        controller_log, deployments)

    details = dict(
        cpu_graph=cpu_image,
        memory_graph=memory_image,
        network_graph=network_image,
        mongo_graph=mongo_image,
        deployments=deployments,
        log_message_chunks=log_message_chunks
    )

    create_html_report(results_dir, details)


def get_log_message_in_timed_chunks(log_file, deployments):
    """Breakdown log into timechunks based on event timeranges in 'deployments'

    """
    deploy_timings = [d.timings for d in deployments['deploys']]

    bootstrap = deployments['bootstrap']
    cleanup = deployments['cleanup']
    all_event_timings = [bootstrap] + deploy_timings + [cleanup]

    raw_details = breakdown_log_by_timeframes(log_file, all_event_timings)

    bs_name = _render_ds_string(bootstrap.start, bootstrap.end)
    cleanup_name = _render_ds_string(cleanup.start, cleanup.end)

    name_lookup = {
        bs_name: 'Bootstrap',
        cleanup_name: 'Kill-Controller',
    }
    for dep in deployments['deploys']:
        name_range = _render_ds_string(
            dep.timings.start, dep.timings.end)
        name_lookup[name_range] = dep.name

    event_details = defaultdict(defaultdict)
    # Outer-layer (i.e. event)
    for event_range in raw_details.keys():
        event_details[event_range]['name'] = name_lookup[event_range]
        event_details[event_range]['logs'] = []

        for log_range in raw_details[event_range].keys():
            timeframe = log_range
            message = '<br/>'.join(raw_details[event_range][log_range])
            event_details[event_range]['logs'].append(
                dict(
                    timeframe=timeframe,
                    message=message))

    return event_details


def create_html_report(results_dir, details):
    # render the html file to the results dir
    with open('./perf_report_template.html', 'rt') as f:
        template = Template(f.read())

    results_output = os.path.join(results_dir, 'report.html')
    with open(results_output, 'wt') as f:
        f.write(template.render(details))


def generate_graph_image(base_dir, results_dir, name, generator):
    metric_files_dir = os.path.join(os.path.abspath(base_dir), results_dir)
    return create_report_graph(metric_files_dir, base_dir, name, generator)


def create_report_graph(rrd_dir, output_dir, name, generator):
    any_file = os.listdir(rrd_dir)[0]
    start, end = get_duration_points(os.path.join(rrd_dir, any_file))
    output_file = os.path.join(
        os.path.abspath(output_dir), '{}.png'.format(name))
    generator(start, end, rrd_dir, output_file)
    print('Created: {}'.format(output_file))
    return output_file


def generate_cpu_graph_image(results_dir):
    return generate_graph_image(
        results_dir, 'aggregation-cpu-average', 'cpu', perf_graphing.cpu_graph)


def generate_memory_graph_image(results_dir):
    return generate_graph_image(
        results_dir, 'memory', 'memory', perf_graphing.memory_graph)


def generate_network_graph_image(results_dir):
    return generate_graph_image(
        results_dir, 'interface-eth0', 'network', perf_graphing.network_graph)


def generate_mongo_graph_image(results_dir):
    dest_path = os.path.join(results_dir, 'mongodb')

    if not create_mongodb_rrd_file(results_dir, dest_path):
        log.error('Failed to create the MongoDB RRD file.')
        return None
    return generate_graph_image(
        results_dir, 'mongodb', 'mongodb', perf_graphing.mongdb_graph)


def create_mongodb_rrd_file(results_dir, destination_dir):
    os.mkdir(destination_dir)
    source_file = os.path.join(results_dir, 'mongodb-stats.log')

    if not os.path.exists(source_file):
        log.warning(
            'Not creating mongodb rrd. Source file not found ({})'.format(
                source_file))
        return False

    dest_file = os.path.join(destination_dir, 'mongodb.rrd')

    first_ts, last_ts, all_data = get_mongodb_stat_data(source_file)
    rrdtool.create(
        dest_file,
        '--start', '{}-10'.format(first_ts),
        '--step', '5',
        'DS:insert:GAUGE:600:0:1000',
        'DS:query:GAUGE:600:0:1000',
        'DS:update:GAUGE:600:0:1000',
        'DS:delete:GAUGE:600:0:1000',
        'RRA:MIN:0.5:1:1200',
        'RRA:MAX:0.5:1:1200',
        'RRA:AVERAGE:0.5:1:120'
    )
    for entry in all_data:
        update_details = '{time}:{i}:{q}:{u}:{d}'.format(
            time=entry[0],
            i=int(entry[1]),
            q=int(entry[2]),
            u=int(entry[3]),
            d=int(entry[4]))
        rrdtool.update(dest_file, update_details)
    return True


def get_mongodb_stat_data(log_file):
    data_lines = []
    with open(log_file, 'rt') as f:
        for line in f:
            details = line.split()
            raw_time = details[-1]
            epoch = int(
                time.mktime(
                    time.strptime(raw_time, '%Y-%m-%dT%H:%M:%SZ')))
            data_lines.append((
                epoch,
                int(details[0].replace('*', '')),
                int(details[1].replace('*', '')),
                int(details[2].replace('*', '')),
                int(details[3].replace('*', ''))))
    return data_lines[0][0], data_lines[-1][0], data_lines


def get_duration_points(rrd_file):
    start = rrdtool.first(rrd_file)
    end = rrdtool.last(rrd_file)

    command = [
        'rrdtool', 'fetch', rrd_file, 'AVERAGE',
        '--start', str(start), '--end', str(end)]
    output = subprocess.check_output(command)

    actual_start = find_actual_start(output)

    return actual_start, end


def find_actual_start(fetch_output):
    # Gets a start timestamp this isn't a NaN.
    for line in fetch_output.splitlines():
        try:
            timestamp, value = line.split(':', 1)
            if not value.startswith(' -nan'):
                return timestamp
        except ValueError:
            pass


def assess_deployment_perf(client, bundle_name='cs:ubuntu'):
    # This is where multiple services are started either at the same time
    # or one after the other etc.
    deploy_start = datetime.utcnow()

    # We possibly want 2 timing details here, one for started (i.e. agents
    # ready) and the other for the workloads to be complete.
    client.deploy(bundle_name)
    client.wait_for_started()
    client.wait_for_workloads()

    deploy_end = datetime.utcnow()
    deploy_timing = TimingData(deploy_start, deploy_end)

    client_details = get_client_details(client)

    return DeployDetails(bundle_name, client_details, deploy_timing)


def get_client_details(client):
    status = client.get_status()
    units = dict()
    for name in status.get_applications().keys():
        units[name] = status.get_service_unit_count(name)
    return units


def setup_system_monitoring(admin_client):
    # Using ssh get into the machine-0 (or all api/state servers)
    # Install the required packages and start up logging of systems collections
    # and mongodb details.

    installer_script_path = _get_static_script_path(SETUP_SCRIPT_PATH)
    collectd_config_path = _get_static_script_path(COLLECTD_CONFIG_PATH)
    installer_script_dest_path = '/tmp/installer.sh'
    runner_script_dest_path = '/tmp/runner.sh'
    collectd_config_dest_file = '/tmp/collectd.config'

    admin_client.juju(
        'scp',
        (collectd_config_path, '0:{}'.format(collectd_config_dest_file)))

    admin_client.juju(
        'scp',
        (installer_script_path, '0:{}'.format(installer_script_dest_path)))
    admin_client.juju('ssh', ('0', 'chmod +x {}'.format(
        installer_script_dest_path)))
    admin_client.juju(
        'ssh',
        ('0', '{installer} {config_file} {output_file}'.format(
            installer=installer_script_dest_path,
            config_file=collectd_config_dest_file,
            output_file=runner_script_dest_path)))

    # Start collection
    # Respawn incase the initial execution fails for whatever reason.
    admin_client.juju('ssh', ('0', '--', 'daemon --respawn {}'.format(
        runner_script_dest_path)))


def _get_static_script_path(script_path):
    full_path = os.path.abspath(__file__)
    current_dir = os.path.dirname(full_path)
    return os.path.join(current_dir, script_path)


def parse_args(argv):
    """Parse all arguments."""
    parser = argparse.ArgumentParser(description="Simple perf/scale testing")
    add_basic_testing_arguments(parser)
    return parser.parse_args(argv)


def main(argv=None):
    args = parse_args(argv)
    configure_logging(args.verbose)
    bs_manager = BootstrapManager.from_args(args)
    assess_perf_test_simple(bs_manager, args.upload_tools)

    return 0


if __name__ == '__main__':
    sys.exit(main())
